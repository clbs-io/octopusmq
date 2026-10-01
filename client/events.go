package client

import "sync/atomic"

// ConnectionLost describes the loss of the stream of a QueueClient or a
// StorageClient: it broke, or its broker stopped leading the cluster.
type ConnectionLost struct {
	// Queue names the queue of a QueueClient, Storage the storage of a
	// StorageClient; the other is empty.
	Queue   string
	Storage string
	// Err is why the stream was lost; it matches ErrStreamBroken or
	// ErrLeaderSwitch.
	Err error
	// LeasesLost reports that the stream held items pulled from the queue, or
	// keys locked in the storage, that were neither settled nor released: the
	// broker freed them, and another consumer may take them.
	LeasesLost bool
}

// lostHandler is the function a Client calls on each lost stream.
type lostHandler struct {
	f atomic.Pointer[func(ConnectionLost)]
}

func (h *lostHandler) call(ev ConnectionLost) {
	if h == nil {
		return
	}
	if f := h.f.Load(); f != nil {
		(*f)(ev)
	}
}

// OnConnectionLost sets f to be called once for each stream of a QueueClient
// or StorageClient of c that is lost, before the operation that noticed the
// loss returns; nil removes it. It applies to the clients c opens after it is
// set, and to those opened before when it replaces an earlier f. f runs on
// that operation's goroutine and must not call the client back; it is meant
// to log, count or signal.
func (c *Client) OnConnectionLost(f func(ConnectionLost)) {
	if c.lost == nil {
		c.lost = &lostHandler{}
	}
	if f == nil {
		c.lost.f.Store(nil)
		return
	}
	c.lost.f.Store(&f)
}
