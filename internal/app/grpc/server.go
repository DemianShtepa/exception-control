package grpc

import (
	"fmt"
	"github.com/DemianShtepa/exception-control/internal/grpc/auth"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type Server struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, port int) *Server {
	gRPCServer := grpc.NewServer()

	auth.Register(gRPCServer)

	return &Server{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *Server) MustRun() {
	err := a.run()
	if err != nil {
		panic(err)
	}
}

func (a *Server) run() error {
	a.log.Info("starting gRPC server")

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("failed to start gRPC server: %w", err)
	}

	a.log.Info("gRPC server is running")

	if err := a.gRPCServer.Serve(listener); err != nil {
		return fmt.Errorf("faied to serve gRPC server: %w", err)
	}

	return nil
}

func (a *Server) Stop() {
	a.log.Info("stopping gRPC server")

	a.gRPCServer.GracefulStop()
}
