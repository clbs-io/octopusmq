#!/usr/bin/env bash


mkdir -p ./pkg/grpcpb
protoc -I=./api/protobuf --go_out=./pkg/grpcpb --go_opt=paths=source_relative --go-grpc_out=./pkg/grpcpb --go-grpc_opt=paths=source_relative ./api/protobuf/queues.proto

mkdir -p ./pkg/grpcmgmtpb
protoc -I=./api/protobuf --go_out=./pkg/grpcmgmtpb --go_opt=paths=source_relative --go-grpc_out=./pkg/grpcmgmtpb --go-grpc_opt=paths=source_relative ./api/protobuf/queues_management.proto

mkdir -p ./pkg/grpcraftcmdpb
protoc -I=./api/protobuf --go_out=./pkg/grpcraftcmdpb --go_opt=paths=source_relative --go-grpc_out=./pkg/grpcraftcmdpb --go-grpc_opt=paths=source_relative ./api/protobuf/raftcommands.proto

mkdir -p ./pkg/grpcclustermgmtpb
protoc -I=./api/protobuf --go_out=./pkg/grpcclustermgmtpb --go_opt=paths=source_relative --go-grpc_out=./pkg/grpcclustermgmtpb --go-grpc_opt=paths=source_relative ./api/protobuf/cluster_management.proto

mkdir -p ./pkg/grpcstoragepb
protoc -I=./api/protobuf --go_out=./pkg/grpcstoragepb --go_opt=paths=source_relative --go-grpc_out=./pkg/grpcstoragepb --go-grpc_opt=paths=source_relative ./api/protobuf/storage.proto

mkdir -p ./pkg/grpcstoragemgmtpb
protoc -I=./api/protobuf --go_out=./pkg/grpcstoragemgmtpb --go_opt=paths=source_relative --go-grpc_out=./pkg/grpcstoragemgmtpb --go-grpc_opt=paths=source_relative ./api/protobuf/storage_management.proto
