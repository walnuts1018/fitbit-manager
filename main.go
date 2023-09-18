package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/walnuts1018/fitbit-manager/config"
	"github.com/walnuts1018/fitbit-manager/handler"
	"github.com/walnuts1018/fitbit-manager/infra/oauth2"
	"github.com/walnuts1018/fitbit-manager/infra/psql"
	"github.com/walnuts1018/fitbit-manager/usecase"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		slog.Error("Error loading config: %v", "error", err)
		os.Exit(1)
	}
	psqlClient, err := psql.NewPSQLClient()
	if err != nil {
		slog.Error("failed to connect to psql", "error", err)
		os.Exit(1)
	}
	defer psqlClient.Close()

	oauth2Client := oauth2.NewOauth2Client()

	tokenUsecase := usecase.NewTokenUsecase(oauth2Client, psqlClient)

	handler, err := handler.NewHandler(tokenUsecase)
	if err != nil {
		slog.Error("Error loading handler: %v", "error", err)
		os.Exit(1)
	}
	err = handler.Run(fmt.Sprintf(":%v", config.ServerPort))
	if err != nil {
		slog.Error("failed to run handler", "error", err)
		os.Exit(1)
	}
}
