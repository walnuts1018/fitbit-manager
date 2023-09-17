package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/walnuts1018/fitbit-manager/config"
	"github.com/walnuts1018/fitbit-manager/handler"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		slog.Error("Error loading config: %v", "error", err)
		os.Exit(1)
	}

	handler, err := handler.NewHandler()
	if err != nil {
		slog.Error("Error loading handler: %v", "error", err)
		os.Exit(1)
	}
	err = handler.Run(fmt.Sprintf(":%v", config.Port))
	if err != nil {
		slog.Error("failed to run handler", "error", err)
		os.Exit(1)
	}
}
