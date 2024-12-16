package fitbit

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/walnuts1018/fitbit-manager/domain"
	"golang.org/x/oauth2"
)

func (c *client) GetName(ctx context.Context, token domain.OAuth2Token) (string, domain.OAuth2Token, error) {
	tokenSource := c.oauth2.TokenSource(ctx, &oauth2.Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    "bearer",
		Expiry:       token.Expiry.StdTime(),
	})

	httpClient := oauth2.NewClient(ctx, tokenSource)
	resp, err := httpClient.Get("https://api.fitbit.com/1/user/-/profile.json")
	if err != nil {
		return "", domain.OAuth2Token{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", domain.OAuth2Token{}, err
	}

	slog.Info("GetName", slog.String("body", string(body)))

	newtoken, err := tokenSource.Token()
	if err != nil {
		return "", domain.OAuth2Token{}, err
	}

	return "", domain.OAuth2Token{
		AccessToken:  newtoken.AccessToken,
		RefreshToken: newtoken.RefreshToken,
		Expiry:       synchro.In[tz.AsiaTokyo](newtoken.Expiry),
	}, nil
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

func (c *client) GetHeartIntraday(ctx context.Context, token domain.OAuth2Token, timeRange domain.FitbitTimeRange, detail domain.HeartDetail) ([]domain.HeartData, error) {
	tokenSource := c.oauth2.TokenSource(ctx, &oauth2.Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    "bearer",
		Expiry:       token.Expiry.StdTime(),
	})
	httpClient := oauth2.NewClient(ctx, tokenSource)
	endpoint := fmt.Sprintf("https://api.fitbit.com/1/user/-/activities/heart/date/%v/1d/%v/time/%v/%v.json", timeRange.Date(), detail, timeRange.StartTime(), timeRange.EndTime())
	resp, err := httpClient.Get(endpoint)
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
