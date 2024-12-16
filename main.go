package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/walnuts1018/fitbit-manager/config"
	"github.com/walnuts1018/fitbit-manager/logger"
	"github.com/walnuts1018/fitbit-manager/tracer"
	"github.com/walnuts1018/fitbit-manager/wire"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	logger.CreateAndSetLogger(cfg.LogLevel, cfg.LogType)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, os.Kill)
	defer stop()

	close, err := tracer.NewTracerProvider(ctx)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to create tracer provider: %v", err))
	}
	defer close()

	router, err := wire.CreateRouter(ctx, cfg)
	if err != nil {
		slog.Error("Failed to create router", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("Server is running", slog.String("port", cfg.ServerPort))

	if err := router.Run(fmt.Sprintf(":%s", cfg.ServerPort)); err != nil {
		slog.Error("Failed to run server", slog.Any("error", err))
		os.Exit(1)
	}
}
