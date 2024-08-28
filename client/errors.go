package client

import (
	"errors"
)

var (
	ErrQueueNotFound      = errors.New("queue not found")
	ErrQueueAlreadyExists = errors.New("queue already exists")
	ErrQueuePaused        = errors.New("queue is paused")
)
