package mockFitbit

import (
	"context"
	"math"
	"time"

	"github.com/walnuts1018/fitbit-manager/domain"
	"github.com/walnuts1018/fitbit-manager/infra/timeJST"
)

type mockOauth2Client struct {
	fitbitClient struct{}
}

func NewMockOauth2Client() domain.FitbitClient {
	return &mockOauth2Client{}
}

func (c *mockOauth2Client) Auth(state string) string {
	return "mockRedirectURL"
}

func (c *mockOauth2Client) Callback(ctx context.Context, code string) (domain.OAuth2Token, error) {
	return domain.OAuth2Token{
		AccessToken:  "mockAccessToken",
		RefreshToken: "mockRefreshToken",
		Expiry:       time.Unix(int64(math.MaxInt64), 0),
		CreatedAt:    timeJST.Now(),
		UpdatedAt:    timeJST.Now(),
	}, nil
}

func (c *mockOauth2Client) NewFitbitClient(ctx context.Context, tokenStore domain.TokenStore) error {
	c.fitbitClient = struct{}{}
	return nil
}

func (c *mockOauth2Client) GetHeartIntraday(date string, startTime string, endTime string, detail domain.HeartDetail) ([]domain.HeartData, error) {
	if startTime == endTime {
		return []domain.HeartData{
			{
				Time:  startTime + ":00",
				Value: 100,
			},
		}, nil
	} else {
		return []domain.HeartData{
			{
				Time:  startTime + ":00",
				Value: 100,
			},
			{
				Time:  endTime + ":00",
				Value: 100,
			},
		}, nil
	}
}
