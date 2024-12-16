package usecase

import (
	"context"
	"fmt"
	"net/url"

	"github.com/walnuts1018/fitbit-manager/util/random"
)

func (u *Usecase) SignIn() (string, string, url.URL, error) {
	state, err := random.String(64, random.AlphanumericSymbols)
	if err != nil {
		return "", "", url.URL{}, fmt.Errorf("failed to generate random string: %w", err)
	}
	redirectURL, verifier, err := u.fitbitClient.GenerateAuthURL(state)
	if err != nil {
		return "", "", url.URL{}, fmt.Errorf("failed to generate auth url: %w", err)
	}
	return state, verifier, redirectURL, nil
}

func (u *Usecase) Callback(ctx context.Context, code string, verifier string) error {
	userID, token, err := u.fitbitClient.Callback(ctx, code, verifier)
	if err != nil {
		return fmt.Errorf("failed to get oauth2 config: %w", err)
	}

	if err := u.tokenStore.SaveOAuth2Token(ctx, userID, token); err != nil {
		return fmt.Errorf("failed to save oauth2 token: %w", err)
	}

	return nil
}
