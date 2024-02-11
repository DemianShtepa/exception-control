package app

import (
	"context"
	"github.com/DemianShtepa/exception-control/internal/app/database"
	grpcapp "github.com/DemianShtepa/exception-control/internal/app/grpc"
	"github.com/DemianShtepa/exception-control/internal/config"
	"log/slog"
)

type App struct {
	GRPCServer *grpcapp.Server
	DB         *database.Database
	log        *slog.Logger
}

func MustInit(ctx context.Context, log *slog.Logger, cfg config.Config) *App {
	db, err := database.New(ctx, log, cfg.DB)
	if err != nil {
		panic(err)
	}
	gRPCServer := grpcapp.New(log, db, cfg.GRPC.Port, cfg.Secret)

	return &App{
		GRPCServer: gRPCServer,
		log:        log,
		DB:         db,
	}
}

func (a *App) MustRun() {
	a.GRPCServer.MustRun()
}

func (a *App) Stop() {
	a.GRPCServer.Stop()
	a.DB.Stop()
}
