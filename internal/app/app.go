package app

import (
	grpcapp "github.com/DemianShtepa/exception-control/internal/app/grpc"
	"log/slog"
)

type App struct {
	GRPCServer *grpcapp.Server
	log        *slog.Logger
	port       int
}

func New(log *slog.Logger, serverPort int) *App {
	gRPCServer := grpcapp.New(log, serverPort)

	return &App{
		GRPCServer: gRPCServer,
		log:        log,
		port:       serverPort,
	}
}
