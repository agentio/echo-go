build:
	go install ./...

all:	rpc grpc connect
	go install ./...

clean:
	go clean
	rm -rf genproto

APIS=$(shell find proto/echo -name "*.proto")

descriptor:
	protoc ${APIS} \
	--proto_path='proto' \
	--include_imports \
	--descriptor_set_out=descriptor.pb

rpc:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	mkdir -p genproto
	protoc ${APIS} \
	--proto_path='proto' \
	--go_opt='module=github.com/agentio/echo-go/genproto' \
	--go_out='genproto'

grpc:
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	mkdir -p genproto
	protoc ${APIS} \
	--proto_path='proto' \
	--go-grpc_opt='module=github.com/agentio/echo-go/genproto' \
	--go-grpc_out='genproto'

connect:
	go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest
	mkdir -p genproto
	protoc ${APIS} \
	--proto_path='proto' \
	--connect-go_opt='module=github.com/agentio/echo-go/genproto' \
	--connect-go_out='genproto'
