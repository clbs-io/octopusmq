package client

import (
	"errors"
)

var (
	ErrQueueNotFound      = errors.New("queue not found")
	ErrQueueAlreadyExists = errors.New("queue already exists")
	ErrQueuePaused        = errors.New("queue is paused")
	ErrQueueClientClosed  = errors.New("queue client closed")
	ErrQueueTimeout       = errors.New("queue operation timeout")

	ErrStorageNotFound      = errors.New("storage not found")
	ErrStorageAlreadyExists = errors.New("storage already exists")
	ErrStorageClientClosed  = errors.New("storage client closed")
	ErrStorageTimeout       = errors.New("storage operation timeout")
	ErrStorageItemNotFound  = errors.New("storage item not found")
)
