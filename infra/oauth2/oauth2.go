package oauth2

import (
	"context"

	"github.com/walnuts1018/fitbit-manager/config"
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

func newOAuth2Conf() *oauth2.Config {
	config := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Endpoint:     oauth2.Endpoint{AuthURL: AuthEndpoint, TokenURL: TokenEndpoint},
		Scopes:       scopes,
	}
	return config
}

func Auth(state string) string {
	conf := newOAuth2Conf()
	url := conf.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return url
}

func Callback(ctx context.Context, code string) (string, error) {
	conf := newOAuth2Conf()
	token, err := conf.Exchange(ctx, code)
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}
