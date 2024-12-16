package usecase

import (
	"context"
	"fmt"

	"github.com/walnuts1018/fitbit-manager/util/random"
)

func (u *Usecase) SignIn() (string, string, error) {
	state, err := random.String(64, random.AlphanumericSymbols)
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
