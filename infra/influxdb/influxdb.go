package influxdb

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/walnuts1018/fitbit-manager/config"
	"github.com/walnuts1018/fitbit-manager/domain"
	"github.com/walnuts1018/fitbit-manager/usecase"
)

type InfluxDBController struct {
	influxdbClient influxdb2.Client
	bucket         string

	queryAPI api.QueryAPI
	writeAPI api.WriteAPIBlocking
}

func NewInfluxDBController(cfg config.InfluxDBConfig) (*InfluxDBController, func()) {
	c := influxdb2.NewClient(cfg.Endpoint.String(), cfg.AuthToken)
	w := c.WriteAPIBlocking(cfg.Org, cfg.Bucket)
	q := c.QueryAPI(cfg.Org)
	return &InfluxDBController{
			influxdbClient: c,
			bucket:         cfg.Bucket,
			queryAPI:       q,
			writeAPI:       w,
		}, func() {
			c.Close()
		}
}

const (
	measurementKey = "heart"
	rateField      = "rate"
	userTag        = "user"
)

func (c *InfluxDBController) RecordHeart(ctx context.Context, userID string, hearts []domain.HeartData) error {
	for _, heart := range hearts {
		p := influxdb2.NewPointWithMeasurement(measurementKey).AddField(rateField, heart.Value).AddTag(userTag, userID).SetTime(heart.Time.StdTime())
		c.writeAPI.WritePoint(context.Background(), p)
	}

	err := c.writeAPI.Flush(ctx)
	if err != nil {
		return fmt.Errorf("failed to flush heart rate: %w", err)
	}
	slog.Info("flushed heart rate datas", "count", len(hearts), "data", hearts)
	return nil
}

func (c *InfluxDBController) GetLatestHeartData(ctx context.Context, userID string) (domain.HeartData, error) {
	h := domain.HeartData{}
	query := fmt.Sprintf(`from(bucket:"%v") |> range(start: 0) |> filter(fn: (r) => r._measurement == "%v") |> filter(fn: (r) => r._field == "%v") |> filter(fn: (r) => r.%v == "%v") |> last()`, c.bucket, measurementKey, rateField, userTag, userID)
	slog.Debug("query", "query", query)
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
		h.Time = synchro.In[tz.AsiaTokyo](result.Record().Time())
	}
	if result.Err() != nil {
		return domain.HeartData{}, fmt.Errorf("query parsing error: %w", result.Err())
	}

	if h == (domain.HeartData{}) {
		return domain.HeartData{}, usecase.ErrNotFound
	}

	return h, nil
}
