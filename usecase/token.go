package usecase

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/walnuts1018/fitbit-manager/domain"
)

type Usecase struct {
	oauth2Client domain.FitbitClient
	tokenStore   domain.TokenStore
	dataStore    domain.DataStore
	heartCache   struct {
		heart     int
		dataAt    time.Time
		UpdatedAt time.Time
	}
}

func NewUsecase(oauth2Client domain.FitbitClient, tokenStore domain.TokenStore, influxdbClient domain.DataStore) *Usecase {
	return &Usecase{
		oauth2Client: oauth2Client,
		tokenStore:   tokenStore,
		dataStore:    influxdbClient,
	}
}

func (u *Usecase) SignIn() (string, string, error) {
	state, err := randStr(64)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate random string: %w", err)
	}
	redirect := u.oauth2Client.Auth(state)
	return state, redirect, nil
}

func (u *Usecase) Callback(ctx context.Context, code string) error {
	cfg, err := u.oauth2Client.Callback(ctx, code)
	if err != nil {
		return fmt.Errorf("failed to get oauth2 config: %w", err)
	}

	err = u.tokenStore.SaveOAuth2Token(cfg)
	if err != nil {
		return fmt.Errorf("failed to save oauth2 token: %w", err)
	}

	return nil
}
func randStr(n int) (string, error) {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	var result string
	for _, v := range b {
		result += string(letters[int(v)%len(letters)])
	}
	return result, nil
}
