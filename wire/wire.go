//go:build wireinject
// +build wireinject

package wire

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/walnuts1018/fitbit-manager/config"
	"github.com/walnuts1018/fitbit-manager/infra/fitbit"
	"github.com/walnuts1018/fitbit-manager/infra/influxdb"
	"github.com/walnuts1018/fitbit-manager/infra/postgres"
	"github.com/walnuts1018/fitbit-manager/router"
	"github.com/walnuts1018/fitbit-manager/router/handler"
	"github.com/walnuts1018/fitbit-manager/usecase"
)

func CreateUsecase(
	ctx context.Context,
	cfg config.Config,
) (*usecase.Usecase, func(), error) {
	wire.Build(
		usecaseSet,
		postgresSet,
		fitbitSet,
		influxdbSet,
		usecase.NewUsecase,
	)
	return &usecase.Usecase{}, nil, nil
}

func CreateRouter(
	ctx context.Context,
	cfg config.Config,
	usecase *usecase.Usecase,
) (*gin.Engine, error) {
	wire.Build(
		routerSet,
		handler.NewHandler,
		router.NewRouter,
	)
	return &gin.Engine{}, nil
}

var usecaseSet = wire.FieldsOf(new(config.Config),
	"ClientID",
	"ClientSecret",
	"PSQLDSN",
	"InfluxDBConfig",
)

var routerSet = wire.FieldsOf(new(config.Config),
	"LogLevel",
	"UserID",
	"CookieSecret",
)

var postgresSet = wire.NewSet(
	postgres.NewPostgres,
	wire.Bind(new(usecase.TokenStore), new(*postgres.PostgresClient)),
)

var fitbitSet = wire.NewSet(
	fitbit.NewFitbitController,
	wire.Bind(new(usecase.FitbitClient), new(*fitbit.FitbitController)),
)

var influxdbSet = wire.NewSet(
	influxdb.NewInfluxDBController,
	wire.Bind(new(usecase.DataStore), new(*influxdb.InfluxDBController)),
)
