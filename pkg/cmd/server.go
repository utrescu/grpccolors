package cmd

import (
	"context"
	"flag"
	"fmt"

	"github.com/utrescu/grpccolors/pkg/protocol/grpc"
	v1 "github.com/utrescu/grpccolors/pkg/service/v1"
)

// Config és la configuració pel servidor
type Config struct {
	GRPCPort string
}

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
	ctx := context.Background()

	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "grpc-port", "", "port en que escoltarà el servidor gRPC")

	flag.Parse()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid port pel servidor gRPC: '%s'", cfg.GRPCPort)
	}

	v1API := v1.NewColorServiceServer()

	return grpc.RunServer(ctx, v1API, cfg.GRPCPort)
}
