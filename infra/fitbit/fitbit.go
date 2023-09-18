package fitbit

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

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
		DataSet []struct {
			Time  string `json:"time"`
			Value int    `json:"value"`
		} `json:"dataset"`
		DatasetInterval int    `json:"datasetInterval"`
		DatasetType     string `json:"datasetType"`
	} `json:"activities-heart-intraday"`
}

func (c *client) GetHeartNow(ctx context.Context) (int, time.Time, error) {
	if c.fclient == nil {
		return 0, time.Time{}, fmt.Errorf("fitbit client is nil")
	}
	now := timeJST.Now()
	if c.heartCache.UpdatedAt.Add(1 * time.Minute).After(now) {
		slog.Info("use cache")
		return c.heartCache.heart, c.heartCache.dataAt, nil
	}

	hourbefore := now.Add(-1 * time.Hour)
	endpoint := fmt.Sprintf("https://api.fitbit.com/1/user/-/activities/heart/date/%v/%v/1sec/time/%v/%v.json", hourbefore.Format("2006-01-02"), now.Format("2006-01-02"), hourbefore.Format("15:04"), now.Format("15:04"))
	resp, err := c.fclient.Get(endpoint)
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("failed to get heart rate: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, time.Time{}, fmt.Errorf("failed to get heart rate: %v", resp.Status)
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var heart heartResult
	err = json.Unmarshal(raw, &heart)
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	data := heart.ActivitiesHeartIntraday.DataSet[len(heart.ActivitiesHeartIntraday.DataSet)-1]
	dtime, err := time.Parse("2006-01-02 15:04:05", now.Format("2006-01-02 ")+data.Time)
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("failed to parse time: %w", err)
	}
	d := now.Sub(dtime)
	if d < (10 * time.Minute) {
		c.heartCache.heart = data.Value
		c.heartCache.dataAt = dtime
		c.heartCache.UpdatedAt = now
		return data.Value, dtime, nil
	} else {
		return 0, time.Time{}, nil
	}
}
