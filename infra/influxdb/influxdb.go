package influxdb

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/walnuts1018/fitbit-manager/config"
	"github.com/walnuts1018/fitbit-manager/domain"
)

type client struct {
	influxdbClient influxdb2.Client
	queryAPI       api.QueryAPI
}

func NewClient(endpoint string) domain.DataStore {
	c := influxdb2.NewClient(endpoint, config.Config.InfluxDBAuthToken)
	q := c.QueryAPI(config.Config.InfluxDBOrg)
	return &client{
		influxdbClient: c,
		queryAPI:       q,
	}
}

func (c *client) Close() {
	c.influxdbClient.Close()
}
