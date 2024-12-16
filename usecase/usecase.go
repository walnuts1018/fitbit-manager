package usecase

import (
	"context"
	"net/url"

	"github.com/walnuts1018/fitbit-manager/domain"
)

type FitbitClient interface {
	GenerateAuthURL(state string) (url.URL, string, error)
	Callback(ctx context.Context, code string, verifier string) (domain.OAuth2Token, error)
	GetName(ctx context.Context, token domain.OAuth2Token) (string, domain.OAuth2Token, error)
	GetHeartData(ctx context.Context, token domain.OAuth2Token, timeRange domain.FitbitTimeRange, detail domain.HeartDetail) ([]domain.HeartData, domain.OAuth2Token, error)
}

type TokenStore interface {
	SaveOAuth2Token(ctx context.Context, userID string, token domain.OAuth2Token) error
	GetOAuth2Token(ctx context.Context, userID string) (domain.OAuth2Token, error)
}

type DataStore interface {
	RecordHeart(ctx context.Context, userID string, hearts []domain.HeartData) error
	GetLatestHeartData(ctx context.Context, userID string) (domain.HeartData, error)
}

type Usecase struct {
	fitbitClient FitbitClient
	tokenStore   TokenStore
	dataStore    DataStore
}

func NewUsecase(fitbitClient FitbitClient, tokenStore TokenStore, dataStore DataStore) *Usecase {
	return &Usecase{
		fitbitClient: fitbitClient,
		tokenStore:   tokenStore,
		dataStore:    dataStore,
	}
}
