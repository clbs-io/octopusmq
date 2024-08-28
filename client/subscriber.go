package client

import (
	"context"

	pb "github.com/clbs-io/octopusmq/api/protobuf"
	"github.com/clbs-io/octopusmq/pkg/grpcpb"
)

// Subscriber is a simple wrapper on top of grpcpb.QueuesService_SubscribeClient,
// since that is internal package and not accessible for outside projects
type Subscriber struct {
	grpcSubscriber grpcpb.QueuesService_SubscribeClient
}

func (sub *Subscriber) Send(m *pb.SubscribeRequest) error {
	return sub.grpcSubscriber.Send(m)
}

func (sub *Subscriber) Recv() (*pb.SubscribeResponse, error) {
	return sub.grpcSubscriber.Recv()
}

func (sub *Subscriber) CloseSend() error {
	return sub.grpcSubscriber.CloseSend()
}

func (sub *Subscriber) Context() context.Context {
	return sub.grpcSubscriber.Context()
}
