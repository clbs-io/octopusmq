package client

import (
	"bytes"
	"testing"

	pb "github.com/clbs-io/octopusmq/api/protobuf"
	"google.golang.org/protobuf/proto"
)

// NewRequestID returns a fresh 16-byte identifier on every call.
func TestNewRequestIDIsUniqueAndSized(t *testing.T) {
	a := NewRequestID()
	b := NewRequestID()

	if len(a) != 16 {
		t.Fatalf("len(a) = %d, want 16", len(a))
	}
	if len(b) != 16 {
		t.Fatalf("len(b) = %d, want 16", len(b))
	}
	if bytes.Equal(a, b) {
		t.Fatalf("two calls to NewRequestID returned the same id: %x", a)
	}
}

// EnqueueRequest.RequestId round-trips through wire encoding unchanged.
func TestEnqueueRequestRoundTripsRequestID(t *testing.T) {
	id := NewRequestID()
	req := &pb.EnqueueRequest{
		Item:      &pb.InputItem{Value: []byte("payload")},
		RequestId: id,
	}

	data, err := proto.Marshal(req)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	got := &pb.EnqueueRequest{}
	if err := proto.Unmarshal(data, got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !bytes.Equal(got.GetRequestId(), id) {
		t.Fatalf("RequestId = %x, want %x", got.GetRequestId(), id)
	}
}

// BatchEnqueueRequest.RequestId round-trips through wire encoding unchanged.
func TestBatchEnqueueRequestRoundTripsRequestID(t *testing.T) {
	id := NewRequestID()
	req := &pb.BatchEnqueueRequest{
		Items:     []*pb.InputItem{{Value: []byte("payload")}},
		RequestId: id,
	}

	data, err := proto.Marshal(req)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	got := &pb.BatchEnqueueRequest{}
	if err := proto.Unmarshal(data, got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !bytes.Equal(got.GetRequestId(), id) {
		t.Fatalf("RequestId = %x, want %x", got.GetRequestId(), id)
	}
}

// GetQueueInfoResponse.FeatureLevel round-trips through wire encoding unchanged.
func TestGetQueueInfoResponseRoundTripsFeatureLevel(t *testing.T) {
	resp := &pb.GetQueueInfoResponse{FeatureLevel: 2}

	data, err := proto.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	got := &pb.GetQueueInfoResponse{}
	if err := proto.Unmarshal(data, got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.GetFeatureLevel() != 2 {
		t.Fatalf("FeatureLevel = %d, want 2", got.GetFeatureLevel())
	}
}
