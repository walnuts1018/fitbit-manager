package logger

import (
	"log/slog"
	"os"

	"github.com/phsym/console-slog"
	"github.com/walnuts1018/fitbit-manager/config"
)

func CreateAndSetLogger(logLevel slog.Level, logType config.LogType) {
	var hander slog.Handler
	switch logType {
	case config.LogTypeText:
		hander = console.NewHandler(os.Stdout, &console.HandlerOptions{
			Level:     logLevel,
			AddSource: logLevel == slog.LevelDebug,
		})
	case config.LogTypeJSON:
		hander = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     logLevel,
			AddSource: logLevel == slog.LevelDebug,
		})
	}

	logger := slog.New(newTraceHandler(hander))
	slog.SetDefault(logger)
}
