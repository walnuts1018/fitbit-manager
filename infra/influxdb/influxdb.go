package influxdb

import (
	"context"
	"fmt"
	"log/slog"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/walnuts1018/fitbit-manager/config"
	"github.com/walnuts1018/fitbit-manager/domain"
	"github.com/walnuts1018/fitbit-manager/infra/timeJST"
)

type client struct {
	influxdbClient influxdb2.Client
	bucket         string

	queryAPI api.QueryAPI
	writeAPI api.WriteAPIBlocking
}

func NewClient(cfg config.InfluxDBConfig) client {
	c := influxdb2.NewClient(cfg.Endpoint, cfg.AuthToken)
	w := c.WriteAPIBlocking(cfg.Org, cfg.Bucket)
	q := c.QueryAPI(cfg.Org)
	return client{
		influxdbClient: c,
		bucket:         cfg.Bucket,
		queryAPI:       q,
		writeAPI:       w,
	}
}

func (c client) Close() {
	c.influxdbClient.Close()
}

func (c client) RecordHeart(ctx context.Context, hearts []domain.HeartData) error {
	for _, heart := range hearts {
		p := influxdb2.NewPointWithMeasurement("heart").AddField("rate", heart.Value).SetTime(*heart.Datatime)
		c.writeAPI.WritePoint(context.Background(), p)
	}

	err := c.writeAPI.Flush(ctx)
	if err != nil {
		return fmt.Errorf("failed to flush heart rate: %w", err)
	}
	slog.Info("flushed heart rate datas", "count", len(hearts), "data", hearts)
	return nil
}

func (c client) GetLastHeartData(ctx context.Context) (domain.HeartData, error) {
	h := domain.HeartData{}
	query := fmt.Sprintf(`from(bucket:"%v") |> range(start: -24h) |> filter(fn: (r) => r._measurement == "heart") |> last()`, c.bucket)
	result, err := c.queryAPI.Query(ctx, query)
	if err != nil {
		return h, fmt.Errorf("failed to query: %w", err)
	}

	for result.Next() {
		v, ok := result.Record().Value().(int64)
		if !ok {
			return domain.HeartData{}, fmt.Errorf("failed to parse value: %v", result.Record().Value())
		}
		h.Value = int(v)
		t := result.Record().Time().In(timeJST.JST)
		h.Datatime = &t
	}
	if result.Err() != nil {
		return domain.HeartData{}, fmt.Errorf("query parsing error: %w", result.Err())
	}
	return h, nil
}
