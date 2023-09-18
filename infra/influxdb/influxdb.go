package influxdb

import (
	"context"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/walnuts1018/fitbit-manager/config"
	"github.com/walnuts1018/fitbit-manager/domain"
)

type client struct {
	influxdbClient influxdb2.Client
	queryAPI       api.QueryAPI
	writeAPI       api.WriteAPIBlocking
}

func NewClient(endpoint string) domain.DataStore {
	c := influxdb2.NewClient(endpoint, config.Config.InfluxDBAuthToken)
	w := c.WriteAPIBlocking(config.Config.InfluxDBAuthToken, config.Config.InfluxDBBucket)
	q := c.QueryAPI(config.Config.InfluxDBOrg)
	return client{
		influxdbClient: c,
		queryAPI:       q,
		writeAPI:       w,
	}
}

func (c client) Close() {
	c.influxdbClient.Close()
}

func (c client) test() {
	p := influxdb2.NewPointWithMeasurement("stat").
		AddTag("unit", "temperature").
		AddField("avg", 23.2).
		AddField("max", 45.0).
		SetTime(time.Now())
	c.writeAPI.WritePoint(context.Background(), p)
}
