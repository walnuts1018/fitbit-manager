package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/walnuts1018/fitbit-manager/config"
	"github.com/walnuts1018/fitbit-manager/handler"
	"github.com/walnuts1018/fitbit-manager/infra/fitbit"
	"github.com/walnuts1018/fitbit-manager/infra/influxdb"
	"github.com/walnuts1018/fitbit-manager/infra/psql"
	"github.com/walnuts1018/fitbit-manager/usecase"
	"golang.org/x/sync/errgroup"
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
	defer influxdbClient.Close()

	usecase := usecase.NewUsecase(oauth2Client, psqlClient, influxdbClient)

	if err := usecase.NewFitbitClient(ctx); err != nil {
		slog.Warn("failed to create fitbit client", "error", err)
	}

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		err := usecase.RecordHeart(ctx)
		if err != nil {
			return fmt.Errorf("failed to record heart: %w", err)
		}

		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return nil
			case <-ticker.C:
				err := usecase.RecordHeart(ctx)
				if err != nil {
					return fmt.Errorf("failed to record heart: %w", err)
				}
			}
		}
	})

	eg.Go(func() error {
		h, err := handler.NewHandler(usecase)
		if err != nil {
			return fmt.Errorf("failed to create handler: %w", err)
		}
		return h.Run(fmt.Sprintf(":%v", config.Config.ServerPort))
	})

	if err := eg.Wait(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
