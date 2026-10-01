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
	// QueueClient or StorageClient opened by a Client repeats such an
	// operation on a new stream while the lost one held nothing, so the error
	// reaches the caller only when the stream held leases (ErrLeasesLost) or no
	// new stream opened within the reopen window.
	ErrStreamBroken = status.Error(codes.Unavailable, "stream broken")

	// ErrLeaderSwitch matches the error of an operation on a stream whose
	// broker no longer leads the cluster: it refused the operation, or its
	// outcome is unknown. A QueueClient or StorageClient opened by a Client
	// treats it as the loss of the stream (see ErrStreamBroken). It is also the
	// error of an operation whose reopen found no leader within the reopen
	// window.
	ErrLeaderSwitch = status.Error(codes.Unavailable, "leadership switching")

	// ErrLeasesLost matches the error of an operation whose stream was lost
	// while it held items pulled from the queue, or keys locked in the
	// storage, that were neither settled nor released: the broker freed them,
	// and another consumer may take them. The operation is not repeated, and
	// whether it took effect is unknown. It also matches the error that lost
	// the stream, ErrStreamBroken or ErrLeaderSwitch.
	ErrLeasesLost = status.Error(codes.Unavailable, "leases lost")

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

// leasesLostError is the error of an operation whose stream was lost while it
// held leases.
type leasesLostError struct {
	cause error
}

func (e *leasesLostError) Error() string {
	return "leases lost: " + e.cause.Error()
}

func (e *leasesLostError) GRPCStatus() *status.Status {
	return status.New(codes.Unavailable, e.Error())
}

func (e *leasesLostError) Is(target error) bool {
	return target == ErrLeasesLost
}

func (e *leasesLostError) Unwrap() error {
	return e.cause
}
