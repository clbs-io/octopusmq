.PHONY: codegen
codegen: generate-protobuf generate-grpc

.PHONY: generate-protobuf
generate-protobuf:
	./hack/generate-protobuf.sh

.PHONY: generate-grpc
generate-grpc:
	./hack/generate-grpc.sh
