# Echo

This repo contains a simple echo service implemented with the official Go [gRPC](https://grpc.io) library ([grpc-go](https://github.com/grpc/grpc-go)) and with the Go implementation of the [Connect RPC protocol](https://connectrpc.com/) ([connect-go](https://github.com/connectrpc/connect-go)).

Servers and clients for all methods are implemented with both stacks and can be used for side-by-side comparisons and examples of how to use both libraries.

Protocol buffer descriptions are in the [proto](/proto) directory.

## Building

`make all` generates necessary code and builds the `echo-go` tool. Code generation requires `protoc`, which can be downloaded [here](https://github.com/protocolbuffers/protobuf/releases). All generated code will be in the `genproto` directory.

`make build` quickly rebuilds when the generated code is already up-to-date.

## Running

`echo-go` is a command-line tool that runs the servers and clients.

## License

Released under the [Apache 2 license](/LICENSE).
