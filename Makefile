.PHONY: codegen
codegen: generate-protobuf generate-grpc

.PHONY: genenerate-protobuf
generate-protobuf:
	( \
		mkdir -p api/protobuf && \
		cd api/protobuf && \
		protoc --go_out=. --go_opt=paths=source_relative  ./messages.proto \
	)

.PHONY: generate-grpc
generate-grpc:
	( \
		mkdir -p pkg/grpcapipb && \
		cd api/protobuf && \
		protoc --proto_path=$(pwd) --go-grpc_out=../../pkg/grpcpb --go-grpc_opt=paths=source_relative  ./queues.proto \
	)
	( \
		mkdir -p pkg/grpcmgmtapipb && \
		cd api/protobuf && \
		protoc --proto_path=$(pwd) --go-grpc_out=../../pkg/grpcmgmtpb --go-grpc_opt=paths=source_relative  ./queues-management.proto \
	)
	( \
		mkdir -p pkg/raftcommandsapipb && \
		cd api/protobuf && \
		protoc --proto_path=$(pwd) --go-grpc_out=../../pkg/grpcraftcmdpb --go-grpc_opt=paths=source_relative  ./raftcommands.proto \
	)
