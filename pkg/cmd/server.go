package cmd

import (
	"context"
	"flag"
	"fmt"

	"github.com/utrescu/grpccolors/pkg/protocol/grpc"
	"github.com/utrescu/grpccolors/pkg/protocol/rest"
	v1 "github.com/utrescu/grpccolors/pkg/service/v1"
)

// Config és la configuració pel servidor
type Config struct {
	GRPCPort string
	HTTPPort string
}

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
	ctx := context.Background()

	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "grpc-port", "", "Port en que escoltarà el servidor gRPC")
	flag.StringVar(&cfg.HTTPPort, "http-port", "", "Port del servei HTTP")
	flag.Parse()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid port pel servidor gRPC: '%s'", cfg.GRPCPort)
	}
	if len(cfg.HTTPPort) == 0 {
		return fmt.Errorf("invalid port pel servidor HTTP: '%s'", cfg.HTTPPort)
	}

	v1API := v1.NewColorServiceServer()

	// Iniciar el gateway
	go func() {
		_ = rest.RunServer(ctx, cfg.GRPCPort, cfg.HTTPPort)
	}()

	return grpc.RunServer(ctx, v1API, cfg.GRPCPort)
}
