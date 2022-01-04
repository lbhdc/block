package handler

import (
	"context"
	"fmt"
	pb "github.com/lbhdc/block/api/v0/net/http"
	"github.com/lbhdc/block/pkg/v0/block"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
)

type HandlerFunc func(ctx context.Context, req *pb.Request) (*pb.Response, error)

type Handler struct {
	pb.UnimplementedHandlerServer
	config  block.HandlerConfig
	handler HandlerFunc
	server  *grpc.Server
}

func NewHandler(name string, handler HandlerFunc) *Handler {
	h := &Handler{
		handler: handler,
		server:  grpc.NewServer(),
	}
	path := os.Getenv("BLOCK_CONFIG")
	cfg := block.NewConfigurationFromFile(path)
	if err := cfg.Valid(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	h.config = cfg.HandlerConfig(name)
	return h
}

func (h *Handler) Handle(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	return h.handler(ctx, req)
}

func (h *Handler) Listen() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", h.config.Port))
	if err != nil {
		log.WithError(err).Error("net.Listen")
		return err
	}
	pb.RegisterHandlerServer(h.server, h)
	return h.server.Serve(lis)
}
