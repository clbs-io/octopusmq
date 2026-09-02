package client

import (
	"errors"
	"io"
	"testing"

	pb "github.com/clbs-io/octopusmq/api/protobuf"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// fakeQueueStream is a bidirectional stream driven by the test.
type fakeQueueStream struct {
	grpc.ClientStream
	sent chan *pb.QueueRequest
	recv chan *pb.QueueResponse
}

func (s *fakeQueueStream) Send(req *pb.QueueRequest) error {
	s.sent <- req
	return nil
}

func (s *fakeQueueStream) Recv() (*pb.QueueResponse, error) {
	r, ok := <-s.recv
	if !ok {
		return nil, io.EOF
	}
	return r, nil
}

func (s *fakeQueueStream) CloseSend() error { return nil }

// fakeStorageStream is a bidirectional stream driven by the test.
type fakeStorageStream struct {
	grpc.ClientStream
	sent chan *pb.StorageRequest
	recv chan *pb.StorageResponse
}

func (s *fakeStorageStream) Send(req *pb.StorageRequest) error {
	s.sent <- req
	return nil
}

func (s *fakeStorageStream) Recv() (*pb.StorageResponse, error) {
	r, ok := <-s.recv
	if !ok {
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

	if err := <-c.errch; !errors.Is(err, io.EOF) {
		t.Fatalf("receiver error = %v, want io.EOF", err)
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

	if err := <-c.errch; !errors.Is(err, io.EOF) {
		t.Fatalf("receiver error = %v, want io.EOF", err)
	}
}

// Callers waiting on a stream that dies get the cause, not a generic error.
func TestQueueClientReportsStreamFailure(t *testing.T) {
	stream := &fakeQueueStream{
		sent: make(chan *pb.QueueRequest, 1),
		recv: make(chan *pb.QueueResponse),
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
	close(stream.recv) // the stream dies with io.EOF

	if err := <-done; !errors.Is(err, io.EOF) {
		t.Fatalf("pending caller error = %v, want io.EOF", err)
	}
	// Later callers learn the same cause instead of blocking or getting a
	// placeholder error.
	if err := c.Noop(); !errors.Is(err, io.EOF) {
		t.Fatalf("subsequent caller error = %v, want io.EOF", err)
	}
}

func TestReduceStreamErr(t *testing.T) {
	other := status.Error(codes.Internal, "boom")
	tests := []struct {
		name string
		err  error
		want error
	}{
		{"canceled", status.Error(codes.Canceled, "canceled"), io.EOF},
		{"deadline exceeded", status.Error(codes.DeadlineExceeded, "deadline"), io.EOF},
		{"other", other, other},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := reducestreamerr(tt.err); got != tt.want {
				t.Fatalf("reducestreamerr(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}
