package grpc

import (
	"fmt"
	"github.com/DemianShtepa/exception-control/internal/app/database"
	"github.com/DemianShtepa/exception-control/internal/app/database/sqlc"
	"github.com/DemianShtepa/exception-control/internal/grpc/auth"
	"github.com/DemianShtepa/exception-control/internal/repository/user/pgsql"
	authservice "github.com/DemianShtepa/exception-control/internal/services/auth"
	"github.com/DemianShtepa/exception-control/internal/services/hash"
	"github.com/DemianShtepa/exception-control/internal/services/token"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type Server struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, db *database.Database, port int, secret string) *Server {
	gRPCServer := grpc.NewServer()

	queries := sqlc.New(db.Pool)
	auth.Register(
		gRPCServer,
		authservice.NewAuth(log, pgsql.NewRepository(db.Pool, queries), &hash.Hasher{}, token.NewGenerator(secret)),
	)

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
