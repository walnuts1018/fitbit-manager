package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
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

	usecase, cleanup, err := wire.CreateUsecase(ctx, cfg)
	if err != nil {
		slog.Error("Failed to create usecase", slog.Any("error", err))
		os.Exit(1)
	}
	defer cleanup()

	from := synchro.Now[tz.AsiaTokyo]().Add(-24 * time.Hour)
	if cfg.RecordStartDatetime != nil {
		from = synchro.Time[tz.AsiaTokyo](*cfg.RecordStartDatetime)
	}

	if err := usecase.RecordHeart(ctx, from, string(cfg.UserID)); err != nil {
		slog.Error("Failed to record heart", slog.Any("error", err))
		os.Exit(1)
	}
	os.Exit(0)
}
