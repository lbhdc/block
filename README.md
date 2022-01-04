# block
block is a pluggable http server, where each handler communicates with the server through a grpc 
interface.

## Build
block depends on `protoc` and `protoc-gen-go` to generate sources before block can be built.
```shell
go generate
go build -o block cmd/block/main.go
```

## Example
See `/examples` for code examples

block uses toml to configure the server and each handler. Each handler is a sub process running 
locally, and needs to be given a unique port to listen to. The entrypoint is the shell command 
needed to start the handler. 
```toml
[server]
addr = "localhost"
port = 8080

[handler.echo]
port = 9090
path = "/echo"
entrypoint = "go run examples/echo/main.go"
```

The example can be run using the following.
```shell
go run cmd/block/main.go -config examples/config.toml

#in another terminal
curl localhost:8080/echo
```