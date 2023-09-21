package mockToken

import "github.com/walnuts1018/fitbit-manager/domain"

type mockTokenClient struct {
	token domain.OAuth2Token
}

func NewMockTokenClient() domain.TokenStore {
	return &mockTokenClient{}
}

func (c *mockTokenClient) SaveOAuth2Token(token domain.OAuth2Token) error {
	c.token = token
	return nil
}

func (c *mockTokenClient) GetOAuth2Token() (domain.OAuth2Token, error) {
	return c.token, nil
}

func (c *mockTokenClient) UpdateOAuth2Token(token domain.OAuth2Token) error {
	c.token = token
	return nil
}

func (c *mockTokenClient) Close() error {
	return nil
}
