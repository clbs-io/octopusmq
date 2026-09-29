package client

import (
	"bytes"
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	pb "github.com/clbs-io/octopusmq/api/protobuf"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// answerFunc answers request req arriving on the stream-th stream the client
// opened, counting from 1; a nil response breaks that stream.
type answerFunc[Req, Resp any] func(stream int, req *Req) *Resp

// scriptedStream is a stream whose responses answerFunc decides.
type scriptedStream[Req, Resp any] struct {
	grpc.ClientStream
	ctx    context.Context
	n      int
	answer answerFunc[Req, Resp]
	recv   chan *Resp
	once   sync.Once
}

func (s *scriptedStream[Req, Resp]) Send(req *Req) error {
	resp := s.answer(s.n, req)
	if resp == nil {
		s.once.Do(func() { close(s.recv) })
		return nil
	}
	s.recv <- resp
	return nil
}

func (s *scriptedStream[Req, Resp]) Recv() (*Resp, error) {
	select {
	case r, ok := <-s.recv:
		if !ok {
			return nil, status.Error(codes.Unavailable, "connection reset by peer")
		}
		return r, nil
	case <-s.ctx.Done():
		return nil, status.FromContextError(s.ctx.Err()).Err()
	}
}

func (s *scriptedStream[Req, Resp]) CloseSend() error { return nil }

// fakeBroker hands out scripted queue and storage streams and counts them.
type fakeBroker struct {
	lock    sync.Mutex
	streams int
	queue   answerFunc[pb.QueueRequest, pb.QueueResponse]
	storage answerFunc[pb.StorageRequest, pb.StorageResponse]
}

func (b *fakeBroker) next() int {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.streams++
	return b.streams
}

func (b *fakeBroker) opened() int {
	b.lock.Lock()
	defer b.lock.Unlock()
	return b.streams
}

func (b *fakeBroker) Connect(ctx context.Context, _ ...grpc.CallOption) (grpc.BidiStreamingClient[pb.QueueRequest, pb.QueueResponse], error) {
	return &scriptedStream[pb.QueueRequest, pb.QueueResponse]{ctx: ctx, n: b.next(), answer: b.queue, recv: make(chan *pb.QueueResponse, 16)}, nil
}

func (b *fakeBroker) StorageConnect(ctx context.Context, _ ...grpc.CallOption) (grpc.BidiStreamingClient[pb.StorageRequest, pb.StorageResponse], error) {
	return &scriptedStream[pb.StorageRequest, pb.StorageResponse]{ctx: ctx, n: b.next(), answer: b.storage, recv: make(chan *pb.StorageResponse, 16)}, nil
}

func (b *fakeBroker) client() *Client {
	return &Client{queueConnect: b, stoConnect: b, logger: zap.NewNop().Sugar()}
}

func queueStatus(req *pb.QueueRequest, code pb.StatusCode) *pb.QueueResponse {
	return &pb.QueueResponse{CorrelationId: req.CorrelationId,
		Response: &pb.QueueResponse_Status{Status: &pb.StatusResponse{Code: code, Message: code.String()}}}
}

func storageStatus(req *pb.StorageRequest, code pb.StatusCode) *pb.StorageResponse {
	return &pb.StorageResponse{CorrelationId: req.CorrelationId,
		Response: &pb.StorageResponse_Status{Status: &pb.StatusResponse{Code: code, Message: code.String()}}}
}

// A leader switch stalls a batch enqueue instead of failing it: the client
// reconnects, the new leader refuses the setup once more while it takes over,
// and the enqueue is then repeated with the request id of the first attempt.
func TestQueueClientRepeatsAnEnqueueAcrossALeaderSwitch(t *testing.T) {
	var ids [][]byte
	b := &fakeBroker{queue: func(stream int, req *pb.QueueRequest) *pb.QueueResponse {
		switch cmd := req.Command.(type) {
		case *pb.QueueRequest_Setup:
			if stream == 2 {
				return queueStatus(req, pb.StatusCode_STATUS_CODE_LEADER_SWITCH)
			}
			return queueStatus(req, pb.StatusCode_STATUS_CODE_OK)
		case *pb.QueueRequest_BatchEnqueue:
			ids = append(ids, cmd.BatchEnqueue.RequestId)
			if stream == 1 {
				return queueStatus(req, pb.StatusCode_STATUS_CODE_LEADER_SWITCH)
			}
			return &pb.QueueResponse{CorrelationId: req.CorrelationId,
				Response: &pb.QueueResponse_BatchEnqueue{BatchEnqueue: &pb.BatchEnqueueResponse{Ids: []uint64{7}}}}
		}
		return queueStatus(req, pb.StatusCode_STATUS_CODE_ERROR)
	}}
	qc, err := b.client().OpenQueue(t.Context(), "q")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer func() { _ = qc.Close() }()

	resp, err := qc.BatchEnqueue(&pb.BatchEnqueueRequest{Items: []*pb.InputItem{{Value: []byte("x")}}})
	if err != nil {
		t.Fatalf("batch enqueue: %v", err)
	}
	if len(resp.Ids) != 1 || resp.Ids[0] != 7 {
		t.Fatalf("ids = %v, want [7]", resp.Ids)
	}
	if got := b.opened(); got != 3 {
		t.Fatalf("streams opened = %d, want 3", got)
	}
	if len(ids) != 2 || len(ids[0]) != 16 || !bytes.Equal(ids[0], ids[1]) {
		t.Fatalf("request ids = %x, want one generated id sent on both attempts", ids)
	}
}

// A stream that breaks under a pull is reopened and the pull repeated.
func TestQueueClientRepeatsAfterABrokenStream(t *testing.T) {
	b := &fakeBroker{queue: func(stream int, req *pb.QueueRequest) *pb.QueueResponse {
		switch req.Command.(type) {
		case *pb.QueueRequest_Setup:
			return queueStatus(req, pb.StatusCode_STATUS_CODE_OK)
		case *pb.QueueRequest_Pull:
			if stream == 1 {
				return nil
			}
			return &pb.QueueResponse{CorrelationId: req.CorrelationId,
				Response: &pb.QueueResponse_Pull{Pull: &pb.PullResponse{Items: []*pb.Item{{Id: 3}}}}}
		}
		return queueStatus(req, pb.StatusCode_STATUS_CODE_ERROR)
	}}
	qc, err := b.client().OpenQueue(t.Context(), "q")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer func() { _ = qc.Close() }()

	resp, err := qc.Pull(&pb.PullRequest{BatchSize: 1})
	if err != nil {
		t.Fatalf("pull: %v", err)
	}
	if len(resp.Items) != 1 || resp.Items[0].Id != 3 {
		t.Fatalf("items = %v, want item 3", resp.Items)
	}
}

// The stall ends with the context the client was opened with, reporting the
// leader switch as Unavailable rather than as an internal error.
func TestQueueClientStopsRepeatingWhenItsContextEnds(t *testing.T) {
	b := &fakeBroker{queue: func(stream int, req *pb.QueueRequest) *pb.QueueResponse {
		if _, ok := req.Command.(*pb.QueueRequest_Setup); ok && stream == 1 {
			return queueStatus(req, pb.StatusCode_STATUS_CODE_OK)
		}
		return queueStatus(req, pb.StatusCode_STATUS_CODE_LEADER_SWITCH)
	}}
	ctx, cancel := context.WithTimeout(t.Context(), 500*time.Millisecond)
	defer cancel()
	qc, err := b.client().OpenQueue(ctx, "q")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer func() { _ = qc.Close() }()

	start := time.Now()
	_, err = qc.Enqueue(&pb.EnqueueRequest{Item: &pb.InputItem{Value: []byte("x")}})
	if !errors.Is(err, ErrLeaderSwitch) || status.Code(err) != codes.Unavailable {
		t.Fatalf("enqueue error = %v, want ErrLeaderSwitch with code Unavailable", err)
	}
	if elapsed := time.Since(start); elapsed < 400*time.Millisecond || elapsed > 2*time.Second {
		t.Fatalf("enqueue gave up after %v, want it to stall until the context ended", elapsed)
	}
	if b.opened() < 2 {
		t.Fatalf("streams opened = %d, want reconnects while stalled", b.opened())
	}
}

// Answers other than a leader switch are returned at once.
func TestQueueClientDoesNotRepeatOtherErrors(t *testing.T) {
	b := &fakeBroker{queue: func(_ int, req *pb.QueueRequest) *pb.QueueResponse {
		if _, ok := req.Command.(*pb.QueueRequest_Setup); ok {
			return queueStatus(req, pb.StatusCode_STATUS_CODE_OK)
		}
		return queueStatus(req, pb.StatusCode_STATUS_CODE_QUEUE_PAUSED)
	}}
	qc, err := b.client().OpenQueue(t.Context(), "q")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer func() { _ = qc.Close() }()

	if _, err := qc.Enqueue(&pb.EnqueueRequest{Item: &pb.InputItem{}}); !errors.Is(err, ErrQueuePaused) {
		t.Fatalf("enqueue error = %v, want ErrQueuePaused", err)
	}
	if got := b.opened(); got != 1 {
		t.Fatalf("streams opened = %d, want 1", got)
	}
}

// A caller-supplied request id is kept.
func TestQueueClientKeepsACallerRequestID(t *testing.T) {
	var got []byte
	b := &fakeBroker{queue: func(_ int, req *pb.QueueRequest) *pb.QueueResponse {
		if cmd, ok := req.Command.(*pb.QueueRequest_Enqueue); ok {
			got = cmd.Enqueue.RequestId
			return &pb.QueueResponse{CorrelationId: req.CorrelationId,
				Response: &pb.QueueResponse_Enqueue{Enqueue: &pb.EnqueueResponse{Id: 1}}}
		}
		return queueStatus(req, pb.StatusCode_STATUS_CODE_OK)
	}}
	qc, err := b.client().OpenQueue(t.Context(), "q")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer func() { _ = qc.Close() }()

	id := []byte("caller")
	if _, err := qc.Enqueue(&pb.EnqueueRequest{Item: &pb.InputItem{}, RequestId: id}); err != nil {
		t.Fatalf("enqueue: %v", err)
	}
	if !bytes.Equal(got, id) {
		t.Fatalf("request id = %q, want %q", got, id)
	}
}

// A storage command refused while leadership moved is repeated, and a key
// listing whose stream broke midway starts over on the new stream.
func TestStorageClientRepeatsAcrossALeaderSwitch(t *testing.T) {
	b := &fakeBroker{storage: func(stream int, req *pb.StorageRequest) *pb.StorageResponse {
		switch req.Command.(type) {
		case *pb.StorageRequest_Setup:
			return storageStatus(req, pb.StatusCode_STATUS_CODE_OK)
		case *pb.StorageRequest_Set:
			if stream == 1 {
				return storageStatus(req, pb.StatusCode_STATUS_CODE_LEADER_SWITCH)
			}
			return storageStatus(req, pb.StatusCode_STATUS_CODE_OK)
		case *pb.StorageRequest_GetKeys:
			return &pb.StorageResponse{CorrelationId: req.CorrelationId,
				Response: &pb.StorageResponse_GetKeysResponse{GetKeysResponse: &pb.StorageGetKeysResponse{Keys: [][]byte{[]byte("a")}, More: true}}}
		case *pb.StorageRequest_GetKeysNext:
			if stream == 2 {
				return nil
			}
			return &pb.StorageResponse{CorrelationId: req.CorrelationId,
				Response: &pb.StorageResponse_GetKeysResponse{GetKeysResponse: &pb.StorageGetKeysResponse{Keys: [][]byte{[]byte("b")}}}}
		}
		return storageStatus(req, pb.StatusCode_STATUS_CODE_ERROR)
	}}
	sc, err := b.client().OpenStorage(t.Context(), "s")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer func() { _ = sc.Close() }()

	if err := sc.Set(&pb.StorageSetRequest{Key: []byte("k")}); err != nil {
		t.Fatalf("set: %v", err)
	}
	keys, err := sc.GetKeys(&pb.StorageGetKeysRequest{})
	if err != nil {
		t.Fatalf("get keys: %v", err)
	}
	if len(keys) != 2 || string(keys[0]) != "a" || string(keys[1]) != "b" {
		t.Fatalf("keys = %q, want [a b] from one complete listing", keys)
	}
	if got := b.opened(); got != 3 {
		t.Fatalf("streams opened = %d, want 3", got)
	}
}
