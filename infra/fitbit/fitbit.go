package fitbit

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

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
	c.tokenStoreCache = tokenStore

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

type heartResult struct {
	ActivitiesHeart []struct {
		CustomHeartRateZones []struct {
			CaloriesOut float64 `json:"caloriesOut"`
			Max         int     `json:"max"`
			Min         int     `json:"min"`
			Minutes     int     `json:"minutes"`
			Name        string  `json:"name"`
		} `json:"customHeartRateZones"`
		DateTime       string `json:"dateTime"`
		HeartRateZones []struct {
			CaloriesOut float64 `json:"caloriesOut"`
			Max         int     `json:"max"`
			Min         int     `json:"min"`
			Minutes     int     `json:"minutes"`
			Name        string  `json:"name"`
		} `json:"heartRateZones"`
		Value string `json:"value"`
	} `json:"activities-heart"`
	ActivitiesHeartIntraday struct {
		DataSet         []domain.HeartData `json:"dataset"`
		DatasetInterval int                `json:"datasetInterval"`
		DatasetType     string             `json:"datasetType"`
	} `json:"activities-heart-intraday"`
}

func (c *client) GetHeartIntraday(ctx context.Context, date string, startTime string, endTime string, detail domain.HeartDetail) ([]domain.HeartData, error) {
	if c.fclient == nil {
		slog.Error("fitbit client is nil, create new client")
		err := c.NewFitbitClient(ctx, c.tokenStoreCache)
		if err != nil {
			return nil, fmt.Errorf("failed to create fitbit client: %w", err)
		}
	}
	endpoint := fmt.Sprintf("https://api.fitbit.com/1/user/-/activities/heart/date/%v/1d/%v/time/%v/%v.json", date, detail, startTime, endTime)
	resp, err := c.fclient.Get(endpoint)

	if err != nil {
		return nil, fmt.Errorf("failed to get heart rate: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get heart rate: %v", resp.Status)
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var heart heartResult
	err = json.Unmarshal(raw, &heart)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	return heart.ActivitiesHeartIntraday.DataSet, nil
}
