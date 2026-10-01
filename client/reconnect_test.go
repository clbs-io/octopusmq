package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
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
	queues  map[int]*scriptedStream[pb.QueueRequest, pb.QueueResponse]
}

// breakQueue breaks the n-th stream while no call waits on it.
func (b *fakeBroker) breakQueue(n int) {
	b.lock.Lock()
	s := b.queues[n]
	b.lock.Unlock()
	s.once.Do(func() { close(s.recv) })
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
	s := &scriptedStream[pb.QueueRequest, pb.QueueResponse]{ctx: ctx, n: b.next(), answer: b.queue, recv: make(chan *pb.QueueResponse, 16)}
	b.lock.Lock()
	if b.queues == nil {
		b.queues = map[int]*scriptedStream[pb.QueueRequest, pb.QueueResponse]{}
	}
	b.queues[s.n] = s
	b.lock.Unlock()
	return s, nil
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

// lostEvents collects the losses a Client reports.
type lostEvents struct {
	lock sync.Mutex
	evs  []ConnectionLost
}

func (l *lostEvents) on(c *Client) *Client {
	c.OnConnectionLost(func(ev ConnectionLost) {
		l.lock.Lock()
		defer l.lock.Unlock()
		l.evs = append(l.evs, ev)
	})
	return c
}

func (l *lostEvents) all() []ConnectionLost {
	l.lock.Lock()
	defer l.lock.Unlock()
	return append([]ConnectionLost(nil), l.evs...)
}

func pullResponse(req *pb.QueueRequest, ids ...uint64) *pb.QueueResponse {
	items := make([]*pb.Item, len(ids))
	for i, id := range ids {
		items[i] = &pb.Item{Id: id}
	}
	return &pb.QueueResponse{CorrelationId: req.CorrelationId, Response: &pb.QueueResponse_Pull{Pull: &pb.PullResponse{Items: items}}}
}

// An enqueue refused while leadership moves, on a stream that holds no item,
// is repeated on a new stream with the request id the client set on the first
// attempt, and the caller sees no error; the loss is reported once.
func TestQueueClientRepeatsAnEnqueueWhileHoldingNothing(t *testing.T) {
	var ids [][]byte
	b := &fakeBroker{queue: func(stream int, req *pb.QueueRequest) *pb.QueueResponse {
		switch cmd := req.Command.(type) {
		case *pb.QueueRequest_Setup:
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
	var lost lostEvents
	qc, err := lost.on(b.client()).OpenQueue(t.Context(), "q")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer func() { _ = qc.Close() }()

	resp, err := qc.BatchEnqueue(&pb.BatchEnqueueRequest{Items: []*pb.InputItem{{Value: []byte("x")}}})
	if err != nil {
		t.Fatalf("batch enqueue: %v", err)
	}
	if len(resp.Ids) != 1 || resp.Ids[0] != 7 || b.opened() != 2 {
		t.Fatalf("ids = %v with %d streams opened, want [7] on a reopened stream", resp.Ids, b.opened())
	}
	if len(ids) != 2 || len(ids[0]) != 16 || !bytes.Equal(ids[0], ids[1]) {
		t.Fatalf("request ids = %x, want one generated id sent on both attempts", ids)
	}
	evs := lost.all()
	if len(evs) != 1 || evs[0].Queue != "q" || evs[0].LeasesLost || !errors.Is(evs[0].Err, ErrLeaderSwitch) {
		t.Fatalf("lost events = %+v, want one leader switch of queue q without lost leases", evs)
	}
}

// A pull whose stream breaks while the stream holds no item is repeated: the
// items the broker may have handed out on it are free again.
func TestQueueClientRepeatsAPullWhileHoldingNothing(t *testing.T) {
	b := &fakeBroker{queue: func(stream int, req *pb.QueueRequest) *pb.QueueResponse {
		switch req.Command.(type) {
		case *pb.QueueRequest_Setup:
			return queueStatus(req, pb.StatusCode_STATUS_CODE_OK)
		case *pb.QueueRequest_Pull:
			if stream == 1 {
				return nil
			}
			return pullResponse(req, 3)
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
	if len(resp.Items) != 1 || resp.Items[0].Id != 3 || b.opened() != 2 {
		t.Fatalf("items = %v with %d streams opened, want item 3 on a reopened stream", resp.Items, b.opened())
	}
}

// A stream lost while it holds pulled items fails the operation with
// ErrLeasesLost, which also names the loss, and reports the lost leases; the
// next operation runs on a new stream, which holds nothing. Items settled
// before the loss do not count.
func TestQueueClientReportsLostLeasesWhileHoldingItems(t *testing.T) {
	for _, settleAll := range []bool{false, true} {
		t.Run(fmt.Sprintf("settled all=%v", settleAll), func(t *testing.T) {
			noops := 0
			b := &fakeBroker{queue: func(stream int, req *pb.QueueRequest) *pb.QueueResponse {
				switch cmd := req.Command.(type) {
				case *pb.QueueRequest_Setup:
					return queueStatus(req, pb.StatusCode_STATUS_CODE_OK)
				case *pb.QueueRequest_Pull:
					return pullResponse(req, 1, 2)
				case *pb.QueueRequest_Commit:
					return &pb.QueueResponse{CorrelationId: req.CorrelationId,
						Response: &pb.QueueResponse_Commit{Commit: &pb.CommitResponse{Ret: int32(len(cmd.Commit.Ids))}}}
				case *pb.QueueRequest_Noop:
					noops++
					if stream == 1 {
						return nil
					}
					return queueStatus(req, pb.StatusCode_STATUS_CODE_OK)
				}
				return queueStatus(req, pb.StatusCode_STATUS_CODE_ERROR)
			}}
			var lost lostEvents
			qc, err := lost.on(b.client()).OpenQueue(t.Context(), "q")
			if err != nil {
				t.Fatalf("open: %v", err)
			}
			defer func() { _ = qc.Close() }()

			if _, err := qc.Pull(&pb.PullRequest{BatchSize: 2}); err != nil {
				t.Fatalf("pull: %v", err)
			}
			settle := []uint64{1}
			if settleAll {
				settle = []uint64{1, 2}
			}
			if _, err := qc.Commit(&pb.CommitRequest{Ids: settle}); err != nil {
				t.Fatalf("commit: %v", err)
			}

			err = qc.Noop()
			evs := lost.all()
			if settleAll {
				if err != nil || noops != 2 {
					t.Fatalf("noop error = %v after %d attempts, want it repeated once the stream held nothing", err, noops)
				}
				if len(evs) != 1 || evs[0].LeasesLost {
					t.Fatalf("lost events = %+v, want one without lost leases", evs)
				}
				return
			}
			if !errors.Is(err, ErrLeasesLost) || !errors.Is(err, ErrStreamBroken) || status.Code(err) != codes.Unavailable {
				t.Fatalf("noop error = %v, want ErrLeasesLost and ErrStreamBroken, Unavailable", err)
			}
			if noops != 1 {
				t.Fatalf("noops sent = %d, want 1: not repeated", noops)
			}
			if len(evs) != 1 || !evs[0].LeasesLost {
				t.Fatalf("lost events = %+v, want one with lost leases", evs)
			}
			if err := qc.Noop(); err != nil {
				t.Fatalf("next noop: %v", err)
			}
			if b.opened() != 2 {
				t.Fatalf("streams opened = %d, want 2", b.opened())
			}
		})
	}
}

// A stream that breaks while no operation waits on it, holding items, fails
// the next operation with ErrLeasesLost, once; the one after runs on a new
// stream.
func TestQueueClientReportsLeasesLostWhileIdle(t *testing.T) {
	b := &fakeBroker{queue: func(_ int, req *pb.QueueRequest) *pb.QueueResponse {
		switch req.Command.(type) {
		case *pb.QueueRequest_Setup, *pb.QueueRequest_Noop:
			return queueStatus(req, pb.StatusCode_STATUS_CODE_OK)
		case *pb.QueueRequest_PullSingle:
			return &pb.QueueResponse{CorrelationId: req.CorrelationId,
				Response: &pb.QueueResponse_PullSingle{PullSingle: &pb.PullSingleResponse{Item: &pb.Item{Id: 4}}}}
		}
		return queueStatus(req, pb.StatusCode_STATUS_CODE_ERROR)
	}}
	var lost lostEvents
	qc, err := lost.on(b.client()).OpenQueue(t.Context(), "q")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer func() { _ = qc.Close() }()

	if _, err := qc.PullSingle(&pb.PullSingleRequest{}); err != nil {
		t.Fatalf("pull: %v", err)
	}
	b.breakQueue(1)
	deadline := time.Now().Add(2 * time.Second)
	for {
		qc.lock.Lock()
		dead := qc.cur.err != nil
		qc.lock.Unlock()
		if dead || time.Now().After(deadline) {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}
	if err := qc.Noop(); !errors.Is(err, ErrLeasesLost) {
		t.Fatalf("noop error = %v, want ErrLeasesLost", err)
	}
	if err := qc.Noop(); err != nil {
		t.Fatalf("next noop: %v", err)
	}
	if evs := lost.all(); len(evs) != 1 || !evs[0].LeasesLost {
		t.Fatalf("lost events = %+v, want one with lost leases", evs)
	}
}

// A reopen the broker refuses while no replica leads yet is retried with a
// backoff, and the operation, as the stream held nothing, is repeated: the
// caller waits for the new leader instead of failing.
func TestQueueClientRetriesARefusedReopen(t *testing.T) {
	var lock sync.Mutex
	noops := map[int]int{}
	b := &fakeBroker{queue: func(stream int, req *pb.QueueRequest) *pb.QueueResponse {
		switch req.Command.(type) {
		case *pb.QueueRequest_Setup:
			if stream == 2 || stream == 3 {
				return queueStatus(req, pb.StatusCode_STATUS_CODE_LEADER_SWITCH)
			}
			return queueStatus(req, pb.StatusCode_STATUS_CODE_OK)
		case *pb.QueueRequest_Noop:
			lock.Lock()
			noops[stream]++
			lock.Unlock()
			if stream == 1 {
				return nil
			}
			return queueStatus(req, pb.StatusCode_STATUS_CODE_OK)
		}
		return queueStatus(req, pb.StatusCode_STATUS_CODE_ERROR)
	}}
	qc, err := b.client().OpenQueue(t.Context(), "q")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer func() { _ = qc.Close() }()

	if err := qc.Noop(); err != nil {
		t.Fatalf("noop: %v, want it to wait out the refused reopens", err)
	}
	lock.Lock()
	defer lock.Unlock()
	if got := b.opened(); got != 4 || noops[1] != 1 || noops[4] != 1 {
		t.Fatalf("streams opened = %d, noops per stream %v; want two refused reopens and the noop repeated once", got, noops)
	}
}

// A reopen that finds no leader within the reopen window the client was opened
// with fails the operation with ErrLeaderSwitch.
func TestQueueClientGivesUpAReopenAfterTheWindow(t *testing.T) {
	b := &fakeBroker{queue: func(stream int, req *pb.QueueRequest) *pb.QueueResponse {
		switch req.Command.(type) {
		case *pb.QueueRequest_Setup:
			if stream > 1 {
				return queueStatus(req, pb.StatusCode_STATUS_CODE_LEADER_SWITCH)
			}
			return queueStatus(req, pb.StatusCode_STATUS_CODE_OK)
		case *pb.QueueRequest_Noop:
			return nil
		}
		return queueStatus(req, pb.StatusCode_STATUS_CODE_ERROR)
	}}
	qc, err := b.client().OpenQueue(t.Context(), "q", WithReopenWindow(500*time.Millisecond))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer func() { _ = qc.Close() }()

	start := time.Now()
	err = qc.Noop()
	if !errors.Is(err, ErrLeaderSwitch) || status.Code(err) != codes.Unavailable {
		t.Fatalf("noop error = %v, want ErrLeaderSwitch with code Unavailable", err)
	}
	if elapsed := time.Since(start); elapsed < 400*time.Millisecond || elapsed > 3*time.Second {
		t.Fatalf("gave up after %v, want about the reopen window", elapsed)
	}
	if b.opened() < 3 {
		t.Fatalf("streams opened = %d, want several reopen attempts", b.opened())
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

// A storage command refused while leadership moved, and a key listing whose
// stream broke midway, are repeated while the stream held no lock: the
// listing starts over. Once the stream holds a lock, a loss fails the command
// with ErrLeasesLost.
func TestStorageClientRepeatsUnlessHoldingLocks(t *testing.T) {
	b := &fakeBroker{storage: func(stream int, req *pb.StorageRequest) *pb.StorageResponse {
		switch req.Command.(type) {
		case *pb.StorageRequest_Setup:
			return storageStatus(req, pb.StatusCode_STATUS_CODE_OK)
		case *pb.StorageRequest_Set:
			if stream == 1 || stream == 3 {
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
		case *pb.StorageRequest_LockAnyWithId:
			return &pb.StorageResponse{CorrelationId: req.CorrelationId,
				Response: &pb.StorageResponse_DataResponse{DataResponse: &pb.StorageDataResponse{Id: 5}}}
		}
		return storageStatus(req, pb.StatusCode_STATUS_CODE_ERROR)
	}}
	var lost lostEvents
	sc, err := lost.on(b.client()).OpenStorage(t.Context(), "s")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer func() { _ = sc.Close() }()

	if err := sc.Set(&pb.StorageSetRequest{Key: []byte("k")}); err != nil {
		t.Fatalf("set: %v, want it repeated", err)
	}
	keys, err := sc.GetKeys(&pb.StorageGetKeysRequest{})
	if err != nil {
		t.Fatalf("get keys: %v", err)
	}
	if len(keys) != 2 || string(keys[0]) != "a" || string(keys[1]) != "b" {
		t.Fatalf("keys = %q, want [a b] from one complete listing", keys)
	}
	if _, err := sc.LockAny(&pb.StorageLockAnyWithIdRequest{}); err != nil {
		t.Fatalf("lock: %v", err)
	}
	if err := sc.Set(&pb.StorageSetRequest{Key: []byte("k")}); !errors.Is(err, ErrLeasesLost) || !errors.Is(err, ErrLeaderSwitch) {
		t.Fatalf("set while holding a lock: %v, want ErrLeasesLost", err)
	}
	if err := sc.Set(&pb.StorageSetRequest{Key: []byte("k")}); err != nil {
		t.Fatalf("next set: %v", err)
	}
	evs := lost.all()
	if len(evs) != 3 || evs[0].LeasesLost || evs[1].LeasesLost || !evs[2].LeasesLost || evs[2].Storage != "s" {
		t.Fatalf("lost events = %+v, want two without and one with lost leases, of storage s", evs)
	}
}

// With a zero reopen window an operation is neither repeated nor waits for a
// reopen: it fails with the loss it met, and the next one reopens the stream
// once.
func TestQueueClientWithoutAReopenWindowNeverRepeats(t *testing.T) {
	noops := 0
	b := &fakeBroker{queue: func(stream int, req *pb.QueueRequest) *pb.QueueResponse {
		switch req.Command.(type) {
		case *pb.QueueRequest_Setup:
			return queueStatus(req, pb.StatusCode_STATUS_CODE_OK)
		case *pb.QueueRequest_Noop:
			noops++
			if stream == 1 {
				return queueStatus(req, pb.StatusCode_STATUS_CODE_LEADER_SWITCH)
			}
			return queueStatus(req, pb.StatusCode_STATUS_CODE_OK)
		}
		return queueStatus(req, pb.StatusCode_STATUS_CODE_ERROR)
	}}
	qc, err := b.client().OpenQueue(t.Context(), "q", WithReopenWindow(0))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer func() { _ = qc.Close() }()

	if err := qc.Noop(); !errors.Is(err, ErrLeaderSwitch) || noops != 1 {
		t.Fatalf("noop error = %v after %d attempts, want ErrLeaderSwitch after one", err, noops)
	}
	if err := qc.Noop(); err != nil || b.opened() != 2 {
		t.Fatalf("next noop error = %v with %d streams opened, want it on a reopened stream", err, b.opened())
	}
}

func TestReopenWindowOfTakesTheLastOption(t *testing.T) {
	if got := reopenWindowOf(nil); got != DefaultReopenWindow {
		t.Fatalf("default window = %v, want %v", got, DefaultReopenWindow)
	}
	opts := []grpc.CallOption{WithReopenWindow(time.Second), grpc.WaitForReady(true), WithReopenWindow(-time.Second)}
	if got := reopenWindowOf(opts); got != 0 {
		t.Fatalf("window = %v, want 0: the last option, negative clamped", got)
	}
}
