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

type Oauth2Client interface {
	Auth(state string) string
	Callback(ctx context.Context, code string) (OAuth2Token, error)
}
