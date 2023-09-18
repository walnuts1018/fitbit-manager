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
	GetHeart(ctx context.Context, startTime time.Time, endTime time.Time) (int, time.Time, error)
}