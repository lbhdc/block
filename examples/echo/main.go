package main

import (
	"context"
	pb "github.com/lbhdc/block/api/v0/net/http"
	sdk "github.com/lbhdc/block/sdk/v0/go/handler"
	"log"
)

func EchoHandler(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	return &pb.Response{Body: []byte(req.Path), Code: 201}, nil
}

func main() {
	handler := sdk.NewHandler("echo", EchoHandler)
	if err := handler.Listen(); err != nil {
		log.Fatalln(err)
	}
}
