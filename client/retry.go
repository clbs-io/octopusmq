package client

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// retryFirst is the wait before the first reconnect; it doubles up to
	// retryMax with every further attempt.
	retryFirst = 100 * time.Millisecond
	retryMax   = 2 * time.Second
)

// leaderSwitchError is a command the broker refused because it does not lead
// the cluster, or answered with an unknown outcome while leadership moved.
type leaderSwitchError struct {
	msg string
}

func (e *leaderSwitchError) Error() string {
	return "leadership switching: " + e.msg
}

func (e *leaderSwitchError) GRPCStatus() *status.Status {
	return status.New(codes.Unavailable, e.Error())
}

func (e *leaderSwitchError) Is(target error) bool {
	return target == ErrLeaderSwitch
}

// retryable reports whether a command that ended with err is repeated on a
// new stream: the stream broke, or the broker answered that leadership is
// moving. Both reach the current leader once the client reconnects.
func retryable(err error) bool {
	return errors.Is(err, ErrStreamBroken) || errors.Is(err, ErrLeaderSwitch)
}

// retrywait waits before the next reconnect, doubling *wait from retryFirst
// up to retryMax. It returns false when ctx ends first.
func retrywait(ctx context.Context, wait *time.Duration) bool {
	*wait = min(max(2**wait, retryFirst), retryMax)
	t := time.NewTimer(*wait)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-t.C:
		return true
	}
}
