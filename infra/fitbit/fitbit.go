package fitbit

import (
	"context"
	"fmt"
	"io"

	"github.com/walnuts1018/fitbit-manager/domain"
	"github.com/walnuts1018/fitbit-manager/infra/timeJST"
	"golang.org/x/oauth2"
)

type tokenSource struct {
	src        oauth2.TokenSource
	tokenStore domain.TokenStore
}

func (s *tokenSource) Token() (*oauth2.Token, error) {
	t, err := s.src.Token()
	if err != nil {
		return nil, err
	}

	token := domain.OAuth2Token{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
		Expiry:       t.Expiry,
		UpdatedAt:    timeJST.Now(),
	}

	err = s.tokenStore.UpdateOAuth2Token(token)
	if err != nil {
		return t, err
	}
	return t, nil
}

func (c *client) NewFitbitClient(ctx context.Context, tokenStore domain.TokenStore) error {
	token, err := tokenStore.GetOAuth2Token()
	if err != nil {
		return fmt.Errorf("failed to get oauth2 token: %w", err)
	}
	oauthToken := &oauth2.Token{
		AccessToken:  token.AccessToken,
		TokenType:    "bearer",
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	}

	oldTokenSource := c.cfg.TokenSource(ctx, oauthToken)
	mySrc := &tokenSource{
		src:        oldTokenSource,
		tokenStore: tokenStore,
	}

	reuseSrc := oauth2.ReuseTokenSource(oauthToken, mySrc)
	c.fclient = oauth2.NewClient(ctx, reuseSrc)
	return nil
}

func (c *client) GetName(ctx context.Context) (string, error) {
	if c.fclient == nil {
		return "", fmt.Errorf("fitbit client is nil")
	}
	resp, err := c.fclient.Get("https://api.fitbit.com/1/user/-/profile.json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	byteArray, _ := io.ReadAll(resp.Body)
	fmt.Println(string(byteArray))
	return "", nil
}
