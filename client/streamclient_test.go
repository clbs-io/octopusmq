package client

import (
	"errors"
	"io"
	"strings"
	"testing"

	pb "github.com/clbs-io/octopusmq/api/protobuf"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// fakeQueueStream is a bidirectional stream driven by the test. Closing recv
// ends the stream with err, or with io.EOF when err is nil.
type fakeQueueStream struct {
	grpc.ClientStream
	sent chan *pb.QueueRequest
	recv chan *pb.QueueResponse
	err  error
}

func (s *fakeQueueStream) Send(req *pb.QueueRequest) error {
	s.sent <- req
	return nil
}

func (s *fakeQueueStream) Recv() (*pb.QueueResponse, error) {
	r, ok := <-s.recv
	if !ok {
		if s.err != nil {
			return nil, s.err
		}
		return nil, io.EOF
	}
	return r, nil
}

func (s *fakeQueueStream) CloseSend() error { return nil }

// fakeStorageStream is a bidirectional stream driven by the test. Closing recv
// ends the stream with err, or with io.EOF when err is nil.
type fakeStorageStream struct {
	grpc.ClientStream
	sent chan *pb.StorageRequest
	recv chan *pb.StorageResponse
	err  error
}

func (s *fakeStorageStream) Send(req *pb.StorageRequest) error {
	s.sent <- req
	return nil
}

func (s *fakeStorageStream) Recv() (*pb.StorageResponse, error) {
	r, ok := <-s.recv
	if !ok {
		if s.err != nil {
			return nil, s.err
		}
		return nil, io.EOF
	}
	return r, nil
}

func (s *fakeStorageStream) CloseSend() error { return nil }

// A response arriving after Close must be discarded instead of being delivered
// to a channel the closing already released.
func TestQueueClientLateResponseAfterClose(t *testing.T) {
	stream := &fakeQueueStream{
		sent: make(chan *pb.QueueRequest, 1),
		recv: make(chan *pb.QueueResponse, 1),
	}
	c := &QueueClient{
		stream:  stream,
		cancel:  func() {},
		logger:  zap.NewNop().Sugar(),
		errch:   make(chan error, 1),
		corrmap: make(map[uint64]chan *pb.QueueResponse),
	}
	go c.receiver()

	done := make(chan error, 1)
	go func() { done <- c.Noop() }()

	req := <-stream.sent
	if err := c.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	if err := <-done; status.Code(err) != codes.Canceled {
		t.Fatalf("pending caller error = %v, want a canceled status", err)
	}

	stream.recv <- &pb.QueueResponse{
		CorrelationId: req.CorrelationId,
		Response: &pb.QueueResponse_Status{
			Status: &pb.StatusResponse{Code: pb.StatusCode_STATUS_CODE_OK},
		},
	}
	close(stream.recv)

	if err := <-c.errch; !errors.Is(err, ErrStreamBroken) {
		t.Fatalf("receiver error = %v, want ErrStreamBroken", err)
	}
	if err := c.Close(); !errors.Is(err, ErrQueueClientClosed) {
		t.Fatalf("second close = %v, want ErrQueueClientClosed", err)
	}
}

// A response arriving after Close must be discarded instead of being delivered
// to a channel the closing already released.
func TestStorageClientLateResponseAfterClose(t *testing.T) {
	stream := &fakeStorageStream{
		sent: make(chan *pb.StorageRequest, 1),
		recv: make(chan *pb.StorageResponse, 1),
	}
	c := &StorageClient{
		stream:  stream,
		cancel:  func() {},
		logger:  zap.NewNop().Sugar(),
		errch:   make(chan error, 1),
		corrmap: make(map[uint64]chan *pb.StorageResponse),
	}
	go c.receiver()

	done := make(chan error, 1)
	go func() { done <- c.Noop() }()

	req := <-stream.sent
	if err := c.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	if err := <-done; status.Code(err) != codes.Canceled {
		t.Fatalf("pending caller error = %v, want a canceled status", err)
	}

	stream.recv <- &pb.StorageResponse{
		CorrelationId: req.CorrelationId,
		Response: &pb.StorageResponse_Status{
			Status: &pb.StatusResponse{Code: pb.StatusCode_STATUS_CODE_OK},
		},
	}
	close(stream.recv)

	if err := <-c.errch; !errors.Is(err, ErrStreamBroken) {
		t.Fatalf("receiver error = %v, want ErrStreamBroken", err)
	}
}

// streamEndings are the ways a stream ends underneath a client. The broker
// losing raft leadership stops its gRPC server, which ends every stream with
// Unavailable.
var streamEndings = []struct {
	name  string
	cause error
}{
	{"server closed the stream", io.EOF},
	{"leader stepped down", status.Error(codes.Unavailable, `closing transport due to: connection error: `+
		`desc = "error reading from server: connection reset by peer", `+
		`received prior goaway: code: NO_ERROR, debug data: "graceful_stop"`)},
	{"deadline exceeded", status.Error(codes.DeadlineExceeded, "context deadline exceeded")},
	{"stream reset", status.Error(codes.Internal, "stream terminated by RST_STREAM")},
}

// checkBrokenStream verifies that err reports a broken stream, names what
// ended it, and cannot be mistaken for the result of the operation itself.
func checkBrokenStream(t *testing.T, who string, err, cause error) {
	t.Helper()
	if !errors.Is(err, ErrStreamBroken) {
		t.Fatalf("%s error = %v, want ErrStreamBroken", who, err)
	}
	if !strings.Contains(err.Error(), cause.Error()) {
		t.Fatalf("%s error = %q, want it to name the cause %q", who, err, cause)
	}
	if code := status.Code(err); code != codes.Unavailable {
		t.Fatalf("%s error code = %v, want Unavailable", who, code)
	}
	if errors.Is(err, ErrQueueTimeout) || errors.Is(err, ErrStorageTimeout) {
		t.Fatalf("%s error = %v reads as an operation timeout", who, err)
	}
}

// Every caller of a queue client whose stream ended, waiting or arriving later,
// learns that the stream is broken and why.
func TestQueueClientReportsBrokenStream(t *testing.T) {
	for _, tt := range streamEndings {
		t.Run(tt.name, func(t *testing.T) {
			stream := &fakeQueueStream{
				sent: make(chan *pb.QueueRequest, 1),
				recv: make(chan *pb.QueueResponse),
				err:  tt.cause,
			}
			c := &QueueClient{
				stream:  stream,
				cancel:  func() {},
				logger:  zap.NewNop().Sugar(),
				errch:   make(chan error, 1),
				corrmap: make(map[uint64]chan *pb.QueueResponse),
			}
			go c.receiver()

			done := make(chan error, 1)
			go func() { done <- c.Noop() }()

			<-stream.sent
			close(stream.recv)

			checkBrokenStream(t, "pending caller", <-done, tt.cause)
			checkBrokenStream(t, "subsequent caller", c.Noop(), tt.cause)
		})
	}
}

// Every caller of a storage client whose stream ended, waiting or arriving
// later, learns that the stream is broken and why.
func TestStorageClientReportsBrokenStream(t *testing.T) {
	for _, tt := range streamEndings {
		t.Run(tt.name, func(t *testing.T) {
			stream := &fakeStorageStream{
				sent: make(chan *pb.StorageRequest, 1),
				recv: make(chan *pb.StorageResponse),
				err:  tt.cause,
			}
			c := &StorageClient{
				stream:  stream,
				cancel:  func() {},
				logger:  zap.NewNop().Sugar(),
				errch:   make(chan error, 1),
				corrmap: make(map[uint64]chan *pb.StorageResponse),
			}
			go c.receiver()

			done := make(chan error, 1)
			go func() { done <- c.Noop() }()

			<-stream.sent
			close(stream.recv)

			checkBrokenStream(t, "pending caller", <-done, tt.cause)
			checkBrokenStream(t, "subsequent caller", c.Noop(), tt.cause)
		})
	}
}
