#!/usr/bin/env bash


mkdir -p ./pkg/grpcpb
protoc -I=./api/protobuf --go_out=./pkg/grpcpb --go_opt=paths=source_relative --go-grpc_out=./pkg/grpcpb --go-grpc_opt=paths=source_relative ./api/protobuf/queues.proto

mkdir -p ./pkg/grpcmgmtpb
protoc -I=./api/protobuf --go_out=./pkg/grpcmgmtpb --go_opt=paths=source_relative --go-grpc_out=./pkg/grpcmgmtpb --go-grpc_opt=paths=source_relative ./api/protobuf/queues_management.proto

mkdir -p ./pkg/grpcraftcmdpb
protoc -I=./api/protobuf --go_out=./pkg/grpcraftcmdpb --go_opt=paths=source_relative --go-grpc_out=./pkg/grpcraftcmdpb --go-grpc_opt=paths=source_relative ./api/protobuf/raftcommands.proto