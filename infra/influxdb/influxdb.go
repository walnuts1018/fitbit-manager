package influxdb

import (
	"context"
	"fmt"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/walnuts1018/fitbit-manager/config"
	"github.com/walnuts1018/fitbit-manager/domain"
)

type client struct {
	influxdbClient influxdb2.Client
	queryAPI       api.QueryAPI
	writeAPI       api.WriteAPIBlocking
}

func NewClient() domain.DataStore {
	c := influxdb2.NewClient(config.Config.InfluxDBEndpoint, config.Config.InfluxDBAuthToken)
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

func (c client) RecordHeart(hearts []domain.HeartData) error {
	var p []*write.Point

	for _, heart := range hearts {
		p = append(p, influxdb2.NewPointWithMeasurement("heart").AddField("rate", heart.Value).SetTime(heart.Datatime))
	}

	err := c.writeAPI.WritePoint(context.Background(), p...)
	if err != nil {
		return fmt.Errorf("failed to write heart rate: %w", err)
	}
	return nil
}
