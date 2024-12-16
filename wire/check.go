package wire

import (
	"github.com/walnuts1018/fitbit-manager/infra/fitbit"
	"github.com/walnuts1018/fitbit-manager/infra/influxdb"
	"github.com/walnuts1018/fitbit-manager/infra/postgres"
	"github.com/walnuts1018/fitbit-manager/usecase"
)

var _ usecase.FitbitClient = &fitbit.FitbitController{}
var _ usecase.TokenStore = &postgres.PostgresClient{}
var _ usecase.DataStore = &influxdb.InfluxDBController{}
