.PHONY: codegen
codegen: generate-protobuf generate-grpc

.PHONY: genenerate-protobuf
generate-protobuf:
	./hack/generate-protobuf.sh

.PHONY: generate-grpc
generate-grpc:
	./hack/generate-grpc.sh
