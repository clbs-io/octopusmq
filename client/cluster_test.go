package client

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/clbs-io/octopusmq/pkg/grpcclustermgmtpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// fakeCluster answers DestroyMember calls from a script, one error per call.
type fakeCluster struct {
	destroy []error
	calls   int
	members []*grpcclustermgmtpb.Member
	listErr error
}

func (f *fakeCluster) ListMembers(context.Context, *emptypb.Empty, ...grpc.CallOption) (*grpcclustermgmtpb.ListMembersResponse, error) {
	if f.listErr != nil {
		return nil, f.listErr
	}
	return &grpcclustermgmtpb.ListMembersResponse{Members: f.members}, nil
}

func (f *fakeCluster) DestroyMember(_ context.Context, req *grpcclustermgmtpb.DestroyMemberRequest, _ ...grpc.CallOption) (*grpcclustermgmtpb.DestroyMemberResponse, error) {
	err := f.destroy[min(f.calls, len(f.destroy)-1)]
	f.calls++
	if err != nil {
		return nil, err
	}
	return &grpcclustermgmtpb.DestroyMemberResponse{RemovedId: uint64(req.Ordinal) + 100}, nil
}

func TestDestroyMemberRetriesUnavailable(t *testing.T) {
	f := &fakeCluster{destroy: []error{
		status.Error(codes.Unavailable, "leadership moved; retry"),
		status.Error(codes.Unavailable, "leadership moved; retry"),
		nil,
	}}
	c := &Client{clusterMgmt: f}
	id, err := c.DestroyMember(t.Context(), 2)
	if err != nil {
		t.Fatalf("destroy: %v", err)
	}
	if id != 102 || f.calls != 3 {
		t.Fatalf("id = %d, calls = %d, want 102 and 3", id, f.calls)
	}
}

func TestDestroyMemberStopsRetryingWhenContextEnds(t *testing.T) {
	f := &fakeCluster{destroy: []error{status.Error(codes.Unavailable, "leadership moved; retry")}}
	c := &Client{clusterMgmt: f}
	ctx, cancel := context.WithTimeout(t.Context(), 300*time.Millisecond)
	defer cancel()
	_, err := c.DestroyMember(ctx, 1)
	if status.Code(err) != codes.Unavailable {
		t.Fatalf("err = %v, want Unavailable", err)
	}
	if f.calls < 2 {
		t.Fatalf("calls = %d, want a retry", f.calls)
	}
}

func TestClusterCallsMapUnimplementedToErrNotSupported(t *testing.T) {
	f := &fakeCluster{destroy: []error{status.Error(codes.Unimplemented, "unknown service")},
		listErr: status.Error(codes.Unimplemented, "unknown service")}
	c := &Client{clusterMgmt: f}
	if _, err := c.DestroyMember(t.Context(), 1); !errors.Is(err, ErrNotSupported) {
		t.Fatalf("destroy err = %v, want ErrNotSupported", err)
	}
	if f.calls != 1 {
		t.Fatalf("calls = %d, want 1", f.calls)
	}
	if _, err := c.ListMembers(t.Context()); !errors.Is(err, ErrNotSupported) {
		t.Fatalf("list err = %v, want ErrNotSupported", err)
	}
}

func TestDestroyMemberPassesFailedPreconditionThrough(t *testing.T) {
	want := status.Error(codes.FailedPrecondition, "replica 1 is still resyncing; wait until it is back")
	f := &fakeCluster{destroy: []error{want}}
	c := &Client{clusterMgmt: f}
	_, err := c.DestroyMember(t.Context(), 2)
	if err != want {
		t.Fatalf("err = %v, want the broker error unchanged", err)
	}
	if f.calls != 1 {
		t.Fatalf("calls = %d, want 1", f.calls)
	}
}

func TestListMembersReturnsMembers(t *testing.T) {
	f := &fakeCluster{members: []*grpcclustermgmtpb.Member{{Ordinal: 0, Id: 1, Voter: true, Leader: true}}}
	c := &Client{clusterMgmt: f}
	got, err := c.ListMembers(t.Context())
	if err != nil || len(got) != 1 || got[0].Id != 1 {
		t.Fatalf("members = %v, err = %v", got, err)
	}
}
