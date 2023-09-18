package usecase

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/walnuts1018/fitbit-manager/domain"
)

type TokenUsecase struct {
	oauth2Client domain.Oauth2Client
	tokenStore   domain.TokenStore
}

func NewTokenUsecase(oauth2Client domain.Oauth2Client, tokenStore domain.TokenStore) *TokenUsecase {
	return &TokenUsecase{
		oauth2Client: oauth2Client,
		tokenStore:   tokenStore,
	}
}

func (u TokenUsecase) Callback(ctx context.Context, code string) error {
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

func (u TokenUsecase) SignIn() (string, string, error) {
	state, err := randStr(64)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate random string: %w", err)
	}
	redirect := u.oauth2Client.Auth(state)
	return state, redirect, nil
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
