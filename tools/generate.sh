#!/usr/bin/env bash

# generate language bindings for the api

# validate that protoc and the necessary plugins are installed
if ! command -v protoc &> /dev/null; then
  printf "%s" "protoc not installed"
  exit 1
fi
if ! command -v protoc-gen-go &> /dev/null; then
  printf "%s" "protoc-gen-go not installed"
  exit 1
fi

# generate go sources
find api -type f -iname "*.proto" -exec protoc -I. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative {} \;