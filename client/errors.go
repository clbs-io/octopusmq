package client

import (
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

	// ErrStreamBroken matches the error of an operation whose stream ended. A
	// QueueClient or StorageClient opened by a Client reconnects and repeats the
	// operation by itself, so the error reaches the caller only once the context
	// the client was opened with ends first.
	ErrStreamBroken = status.Error(codes.Unavailable, "stream broken")

	// ErrLeaderSwitch matches the error of an operation the broker refused
	// because it no longer leads the cluster. A QueueClient or StorageClient
	// opened by a Client reconnects to the new leader and repeats the operation
	// by itself, so the error reaches the caller only once the context the
	// client was opened with ends first.
	ErrLeaderSwitch = status.Error(codes.Unavailable, "leadership switching")

	// ErrNotSupported matches the error of a call the broker does not implement,
	// such as a cluster management call sent to a broker that predates it.
	ErrNotSupported = status.Error(codes.Unimplemented, "not supported by the broker")
)

// brokenStreamError is what a client reports once its stream has ended. Every
// error Send or Recv returns ends a gRPC stream, so the client wraps each of them
// in it. The message names the error that ended the stream, but the status is
// always Unavailable: a stream that ended on a missed deadline must not read as
// an operation that timed out, and one the broker tore down while losing raft
// leadership must not read as a paused queue.
type brokenStreamError struct {
	cause error
}

func (e *brokenStreamError) Error() string {
	return "stream broken: " + e.cause.Error()
}

func (e *brokenStreamError) GRPCStatus() *status.Status {
	return status.New(codes.Unavailable, e.Error())
}

func (e *brokenStreamError) Is(target error) bool {
	return target == ErrStreamBroken
}
