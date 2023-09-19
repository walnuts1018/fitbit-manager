package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/walnuts1018/fitbit-manager/config"
	"github.com/walnuts1018/fitbit-manager/handler"
	"github.com/walnuts1018/fitbit-manager/infra/fitbit"
	"github.com/walnuts1018/fitbit-manager/infra/influxdb"
	"github.com/walnuts1018/fitbit-manager/infra/psql"
	"github.com/walnuts1018/fitbit-manager/usecase"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		slog.Error("Error loading config: %v", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()

	psqlClient, err := psql.NewPSQLClient()
	if err != nil {
		slog.Error("failed to connect to psql", "error", err)
		os.Exit(1)
	}
	defer psqlClient.Close()

	oauth2Client := fitbit.NewOauth2Client()

	influxdbClient := influxdb.NewClient()

	usecase := usecase.NewUsecase(oauth2Client, psqlClient, influxdbClient)

	err = usecase.NewFitbitClient(ctx)
	if err != nil {
		slog.Warn("failed to create fitbit client", "error", err)
	}

	go func() {
		err := usecase.RecordHeart(ctx)
		if err != nil {
			slog.Error("Failed to record heart", "error", err)
		}
	}()

	handler, err := handler.NewHandler(usecase)
	if err != nil {
		slog.Error("Error loading handler: %v", "error", err)
		os.Exit(1)
	}
	err = handler.Run(fmt.Sprintf(":%v", config.Config.ServerPort))
	if err != nil {
		slog.Error("failed to run handler", "error", err)
		os.Exit(1)
	}
}
