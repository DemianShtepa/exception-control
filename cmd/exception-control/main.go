package main

import (
	"context"
	"github.com/DemianShtepa/exception-control/internal/app"
	"github.com/DemianShtepa/exception-control/internal/config"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	cfg := config.MustLoad()

	logger := setupLogger(cfg.Env)

	logger.Info("starting application", slog.String("env", cfg.Env))

	application := app.MustInit(ctx, logger, *cfg)
	go application.MustRun()

	<-ctx.Done()

	application.Stop()
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
