package client

import (
	"io"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrQueueNotFound      = status.Error(codes.NotFound, "queue not found")
	ErrQueueAlreadyExists = status.Error(codes.AlreadyExists, "queue already exists")
	ErrQueuePaused        = status.Error(codes.Unavailable, "queue is paused")
	ErrQueueClientClosed  = status.Error(codes.Canceled, "queue client closed")
	ErrQueueTimeout       = status.Error(codes.DeadlineExceeded, "queue operation timeout")

	ErrStorageNotFound      = status.Error(codes.NotFound, "storage not found")
	ErrStorageAlreadyExists = status.Error(codes.AlreadyExists, "storage already exists")
	ErrStorageClientClosed  = status.Error(codes.Canceled, "storage client closed")
	ErrStorageTimeout       = status.Error(codes.DeadlineExceeded, "storage operation timeout")
	ErrStorageKeyNotFound   = status.Error(codes.NotFound, "storage key not found")
)

// reducestreamerr maps a terminal transport error of a bidirectional stream onto
// io.EOF, so that callers can treat a broken stream uniformly regardless of how
// it was torn down. Any other error is returned unchanged.
func reducestreamerr(err error) error {
	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.Canceled, codes.DeadlineExceeded:
			return io.EOF
		}
	}
	return err
}
