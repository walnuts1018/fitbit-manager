package usecase

import (
	"context"

	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/walnuts1018/fitbit-manager/domain"
)

type FitbitClient interface {
	Auth(state string) string
	Callback(ctx context.Context, code string) (domain.OAuth2Token, error)
	NewFitbitClient(ctx context.Context, tokenStore TokenStore) error
	GetHeartIntraday(ctx context.Context, date string, startTime string, endTime string, detail domain.HeartDetail) ([]domain.HeartData, error)
}

type TokenStore interface {
	SaveOAuth2Token(token domain.OAuth2Token) error
	GetOAuth2Token() (domain.OAuth2Token, error)
	UpdateOAuth2Token(token domain.OAuth2Token) error
	Close() error
}

type DataStore interface {
	RecordHeart(ctx context.Context, rates []domain.HeartData) error
	GetLastHeartData(ctx context.Context) (domain.HeartData, error)
	Close()
}

type Usecase struct {
	oauth2Client FitbitClient
	tokenStore   TokenStore
	dataStore    DataStore
	heartCache   struct {
		heart     int
		dataAt    synchro.Time[tz.AsiaTokyo]
		UpdatedAt synchro.Time[tz.AsiaTokyo]
	}
}

func NewUsecase(oauth2Client FitbitClient, tokenStore TokenStore, influxdbClient DataStore) *Usecase {
	return &Usecase{
		oauth2Client: oauth2Client,
		tokenStore:   tokenStore,
		dataStore:    influxdbClient,
	}
}
