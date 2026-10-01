package client

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc"
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

// retryable reports whether a command that ended with err leaves its stream
// to be reopened: the stream broke, or the broker answered that leadership is
// moving. A new stream reaches the current leader.
func retryable(err error) bool {
	return errors.Is(err, ErrStreamBroken) || errors.Is(err, ErrLeaderSwitch)
}

// DefaultReopenWindow is how long an operation of a QueueClient or
// StorageClient waits for its stream to reopen and repeats itself while the
// stream holds nothing, unless WithReopenWindow sets another window.
const DefaultReopenWindow = 30 * time.Second

// reopenWindowOption carries the window WithReopenWindow sets.
type reopenWindowOption struct {
	grpc.EmptyCallOption
	window time.Duration
}

// WithReopenWindow sets, among the options of Client.OpenQueue or
// Client.OpenStorage, how long an operation of the client waits for a lost
// stream to reopen, while the broker is unreachable or no replica leads, and
// repeats itself while the stream holds nothing. Zero or less tries the
// reopen once and never repeats an operation. Other calls ignore it.
func WithReopenWindow(window time.Duration) grpc.CallOption {
	return reopenWindowOption{window: max(window, 0)}
}

// reopenWindowOf returns the window the last WithReopenWindow among opts sets,
// or DefaultReopenWindow.
func reopenWindowOf(opts []grpc.CallOption) time.Duration {
	window := DefaultReopenWindow
	for _, o := range opts {
		if w, ok := o.(reopenWindowOption); ok {
			window = w.window
		}
	}
	return window
}

// reopening runs reconnect, which opens a new stream, until it succeeds or
// fails with an error that a later attempt cannot mend. While the broker is
// unreachable, or refuses the setup because no replica leads yet, nothing ran
// and the stream holds nothing, so the attempt is repeated, backing off from
// retryFirst up to retryMax, for at most window or until ctx ends; the last
// error is returned then.
func reopening(ctx context.Context, window time.Duration, reconnect func() error) error {
	ctx, cancel := context.WithTimeout(ctx, window)
	defer cancel()
	var wait time.Duration
	for {
		err := reconnect()
		if err == nil || !retryable(err) || !retrywait(ctx, &wait) {
			return err
		}
	}
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
