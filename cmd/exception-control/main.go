package main

import (
	"github.com/DemianShtepa/exception-control/internal/app"
	"github.com/DemianShtepa/exception-control/internal/config"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	logger := setupLogger(cfg.Env)

	logger.Info("starting application", slog.String("env", cfg.Env))

	application := app.New(logger, cfg.GRPC.Port)
	application.GRPCServer.MustRun()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	var level slog.Level

	switch env {
	case envLocal, envDev:
		level = slog.LevelDebug
	case envProd:
		level = slog.LevelInfo
	default:
		panic("unsupported env provided")
	}

	log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))

	return log
}
