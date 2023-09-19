package fitbit

import (
	"context"
	"net/http"

	"github.com/walnuts1018/fitbit-manager/config"
	"github.com/walnuts1018/fitbit-manager/domain"
	"github.com/walnuts1018/fitbit-manager/infra/timeJST"
	"golang.org/x/oauth2"
)

var (
	AuthEndpoint  = "https://www.fitbit.com/oauth2/authorize"
	TokenEndpoint = "https://api.fitbit.com/oauth2/token"
	scopes        = []string{
		"activity",
		"cardio_fitness",
		"electrocardiogram",
		"heartrate",
		"location",
		"nutrition",
		"oxygen_saturation",
		"profile",
		"respiratory_rate",
		"settings",
		"sleep",
		"social",
		"temperature",
		"weight",
	}
)

type client struct {
	cfg     *oauth2.Config
	fclient *http.Client //fitbit client
}

func NewOauth2Client() domain.FitbitClient {
	return &client{
		cfg: &oauth2.Config{
			ClientID:     config.Config.ClientID,
			ClientSecret: config.Config.ClientSecret,
			Endpoint:     oauth2.Endpoint{AuthURL: AuthEndpoint, TokenURL: TokenEndpoint},
			Scopes:       scopes,
		},
	}
}

func (c *client) Auth(state string) string {
	url := c.cfg.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return url
}

func (c *client) Callback(ctx context.Context, code string) (domain.OAuth2Token, error) {
	token, err := c.cfg.Exchange(ctx, code)
	if err != nil {
		return domain.OAuth2Token{}, err
	}

	cfg := domain.OAuth2Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
		CreatedAt:    timeJST.Now(),
		UpdatedAt:    timeJST.Now(),
	}
	return cfg, nil
}
