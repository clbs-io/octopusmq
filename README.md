# OctopusMQ

**OctopusMQ** is a distributed, highly-available message queue and key-value store system built on the Raft consensus algorithm. This repository contains the public API definitions, Protocol Buffer schemas, and Go client libraries for integrating with OctopusMQ.

## Overview

OctopusMQ provides two core capabilities:

- **Priority Message Queue**: A distributed priority queue with TTL support, key-based grouping, and batch operations
- **Key-Value Store**: A consistent key-value storage system with locking mechanisms and TTL support

Both systems are backed by Raft consensus, ensuring strong consistency and high availability across distributed deployments.

## Features

### Message Queue

- **Priority-based queuing** with up to 32 configurable priority levels
- **TTL (Time-To-Live)** support for automatic message expiration
- **Key-based grouping** for correlated message delivery
- **Batch operations** for pushing and pulling multiple messages
- **Optional compression** for message payloads
- **Requeue capability** for failed message processing
- **Pull-based subscriber pattern** for message consumption

### Key-Value Store

- **Basic operations**: Set, Get, Delete with TTL support
- **Locking mechanism**: Exclusive key locking (LockKey, LockAny)
- **TTL support**: Automatic key expiration
- **Optional compression** for stored values
- **Bulk operations**: Retrieve multiple keys by ID
- **State tracking**: Track active/busy keys during processing

## Installation

```bash
go get github.com/clbs-io/octopusmq@latest
```

## Quick Start

### Message Queue Client

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/clbs-io/octopusmq/client"
    "google.golang.org/protobuf/types/known/durationpb"
)

func main() {
    // Create a new queue client
    qClient, err := client.NewQueueClient("localhost:4123", "my-queue")
    if err != nil {
        log.Fatal(err)
    }
    defer qClient.Close()

    ctx := context.Background()

    // Enqueue a message
    item := &client.InputItem{
        Priority: 10,
        Key:      []byte("order-123"),
        Value:    []byte("process this order"),
        TTL:      durationpb.New(5 * time.Minute),
    }

    id, err := qClient.Enqueue(ctx, item)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Enqueued message with ID: %d", id)

    // Pull messages
    items, err := qClient.Pull(ctx, 10)
    if err != nil {
        log.Fatal(err)
    }

    for _, item := range items {
        log.Printf("Processing message ID %d: %s", item.Id, string(item.Value))

        // Commit when done
        if err := qClient.CommitSingle(ctx, item.Id); err != nil {
            log.Printf("Failed to commit: %v", err)
        }
    }
}
```

### Key-Value Store Client

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/clbs-io/octopusmq/client"
    "google.golang.org/protobuf/types/known/durationpb"
)

func main() {
    // Create a new storage client
    sClient, err := client.NewStorageClient("localhost:4123", "my-storage")
    if err != nil {
        log.Fatal(err)
    }
    defer sClient.Close()

    ctx := context.Background()

    // Set a key-value pair
    id, err := sClient.SetKey(ctx, []byte("user:123"), []byte(`{"name":"John"}`),
        durationpb.New(10*time.Minute), false)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Set key with ID: %d", id)

    // Get a value
    value, err := sClient.GetKey(ctx, []byte("user:123"))
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Retrieved value: %s", string(value))

    // Lock a key
    lockedID, err := sClient.LockKey(ctx, []byte("user:123"),
        durationpb.New(30*time.Second))
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Locked key with ID: %d", lockedID)

    // Unlock when done
    if err := sClient.UnlockKey(ctx, lockedID); err != nil {
        log.Fatal(err)
    }
}
```

## API Reference

### gRPC Services

OctopusMQ exposes four gRPC service interfaces:

1. **QueuesService** ([queues.proto](api/protobuf/queues.proto)): Message queue operations
2. **StorageService** ([storage.proto](api/protobuf/storage.proto)): Key-value storage operations
3. **ManagementService** ([queues_management.proto](api/protobuf/queues_management.proto)): Queue management operations
4. **StorageManagementService** ([storage_management.proto](api/protobuf/storage_management.proto)): Storage management operations

### Protocol Buffers

All API definitions are available in the [api/protobuf](api/protobuf/) directory:

- [messages.proto](api/protobuf/messages.proto) - Core message structures and request/response types
- [queues.proto](api/protobuf/queues.proto) - Queue service definitions
- [storage.proto](api/protobuf/storage.proto) - Storage service definitions
- [queues_management.proto](api/protobuf/queues_management.proto) - Queue management service
- [storage_management.proto](api/protobuf/storage_management.proto) - Storage management service

### Go Client Library

The [client](client/) package provides high-level Go clients:

- **QueueClient**: Simplified interface for queue operations
- **StorageClient**: Simplified interface for key-value operations

Both clients handle connection management, correlation IDs, and error handling automatically.

## Queue Operations

### Enqueue Messages

```go
// Single message
id, err := qClient.Enqueue(ctx, item)

// Batch enqueue
ids, err := qClient.BatchEnqueue(ctx, items)
```

### Pull Messages

```go
// Pull up to N messages
items, err := qClient.Pull(ctx, maxMessages)

// Pull single message
item, err := qClient.PullSingle(ctx)
```

### Commit Messages

```go
// Commit single message
err := qClient.CommitSingle(ctx, itemID)

// Commit multiple messages
err := qClient.Commit(ctx, itemIDs)
```

### Requeue Messages

```go
// Requeue with new TTL and priority
err := qClient.RequeueSingle(ctx, itemID, newTTL, newPriority)
```

### Delete Messages

```go
// Delete single message
err := qClient.DeleteSingle(ctx, itemID)

// Delete multiple messages
err := qClient.Delete(ctx, itemIDs)
```

## Storage Operations

### Set Keys

```go
id, err := sClient.SetKey(ctx, key, value, ttl, compress)
```

### Get Keys

```go
// Get single key
value, err := sClient.GetKey(ctx, key)

// Get multiple keys by ID
items, err := sClient.GetKeys(ctx, ids)
```

### Lock Keys

```go
// Lock specific key
id, err := sClient.LockKey(ctx, key, ttl)

// Lock any available key with prefix
id, key, err := sClient.LockAny(ctx, keyPrefix, ttl)
```

### Unlock Keys

```go
err := sClient.UnlockKey(ctx, id)
```

### Delete Keys

```go
err := sClient.DeleteKey(ctx, key)
```

## Connection Management

Both clients support automatic reconnection and leader discovery:

```go
client, err := client.NewQueueClient("server:4123", "queue-name")
// Client automatically handles:
// - Connection establishment
// - Stream management
// - Correlation ID tracking
// - Error handling
```

## Status Codes

The API uses the following status codes:

| Code | Description |
|------|-------------|
| `STATUS_CODE_OK` | Operation successful |
| `STATUS_CODE_ERROR` | General error |
| `STATUS_CODE_TIMEOUT` | Operation timed out |
| `STATUS_CODE_NOT_FOUND` | Resource not found |
| `STATUS_CODE_INVALID_ARGUMENT` | Invalid request parameter |
| `STATUS_CODE_ALREADY_EXISTS` | Resource already exists |
| `STATUS_CODE_QUEUE_NOT_FOUND` | Queue does not exist |
| `STATUS_CODE_ITEM_BUSY` | Item is locked/busy |
| `STATUS_CODE_QUEUE_PAUSED` | Queue is paused |
| `STATUS_CODE_LEADER_SWITCH` | Raft leader changed |
| `STATUS_CODE_STORAGE_NOT_FOUND` | Storage does not exist |
| `STATUS_CODE_ITEM_NOT_BUSY` | Item is not locked |
| `STATUS_CODE_ITEM_NOT_FOUND` | Item does not exist |

## Error Handling

The client library provides typed errors:

```go
import "github.com/clbs-io/octopusmq/client"

items, err := qClient.Pull(ctx, 10)
if err != nil {
    switch {
    case errors.Is(err, client.ErrQueueNotFound):
        // Handle queue not found
    case errors.Is(err, client.ErrTimeout):
        // Handle timeout
    default:
        // Handle other errors
    }
}
```

## Building from Source

### Prerequisites

- Go 1.25 or later
- Protocol Buffers compiler (protoc)
- gRPC Go plugins

### Generate Protocol Buffers

```bash
make generate
```

### Build

```bash
go build ./...
```

### Run Tests

```bash
go test ./...
```

## Deployment

OctopusMQ server is designed to run in Kubernetes as a StatefulSet. Contact dDash s.r.o. for deployment packages and enterprise support.

## License

Copyright © 2024-2025 dDash s.r.o. All rights reserved.

This software is proprietary and confidential. The source code is not publicly available. This repository contains only the public API definitions and client libraries for integration purposes.

## Support

For technical support, enterprise licensing, or deployment assistance, please contact:

**dDash s.r.o.**
- Company: [https://ddash.cz](https://ddash.cz)
- Brand: [https://www.cybroslabs.com](https://www.cybroslabs.com)
- Email: info@cybroslabs.com

## Contributing

This is a proprietary product. The API definitions and client libraries are maintained by dDash s.r.o.

## Acknowledgments

OctopusMQ is built on top of excellent open-source technologies:

- [gRPC](https://grpc.io/) - High-performance RPC framework
- [Protocol Buffers](https://protobuf.dev/) - Language-neutral data serialization
- [HashiCorp Raft](https://github.com/hashicorp/raft) - Consensus algorithm implementation (server-side)
