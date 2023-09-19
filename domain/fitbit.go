package domain

import (
	"context"
	"time"
)

type OAuth2Token struct {
	AccessToken  string
	RefreshToken string
	Expiry       time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type FitbitClient interface {
	Auth(state string) string
	Callback(ctx context.Context, code string) (OAuth2Token, error)
	NewFitbitClient(ctx context.Context, tokenStore TokenStore) error
	GetHeartIntraday(date string, startTime string, endTime string, detail HeartDetail) ([]HeartData, error)
}

type HeartData struct {
	Datatime time.Time
	Time     string `json:"time"`
	Value    int    `json:"value"`
}

type HeartDetail string

const (
	HeartDetailOneSecond     HeartDetail = "1sec"
	HeartDetailOneMinute     HeartDetail = "1min"
	HeartDetailFiveMinute    HeartDetail = "5min"
	HeartDetailFifteenMinute HeartDetail = "15min"
)
