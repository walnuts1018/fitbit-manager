package fitbit

import (
	"context"
	"fmt"
	"net/url"

	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/walnuts1018/fitbit-manager/config"
	"github.com/walnuts1018/fitbit-manager/domain"
	"golang.org/x/oauth2"
)

const (
	AuthEndpoint  = "https://www.fitbit.com/oauth2/authorize"
	TokenEndpoint = "https://api.fitbit.com/oauth2/token"
)

var (
	scopes = []string{
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

type FitbitController struct {
	oauth2 *oauth2.Config
}

func NewFitbitController(serverURL config.ServerURL, clientID config.ClientID, clientSecret config.ClientSecret) (*FitbitController, error) {
	redirectURL, err := url.JoinPath(string(serverURL), "/callback")
	if err != nil {
		return nil, err
	}

	return &FitbitController{
		oauth2: &oauth2.Config{
			ClientID:     string(clientID),
			ClientSecret: string(clientSecret),
			Endpoint:     oauth2.Endpoint{AuthURL: AuthEndpoint, TokenURL: TokenEndpoint},
			Scopes:       scopes,
			RedirectURL:  redirectURL,
		},
	}, nil
}

func (c *FitbitController) GenerateAuthURL(state string) (url.URL, string, error) {
	verifier := oauth2.GenerateVerifier()
	s := c.oauth2.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))

	u, err := url.Parse(s)
	if err != nil {
		return url.URL{}, "", err
	}

	return *u, verifier, nil
}

func (c *FitbitController) Callback(ctx context.Context, code string, verifier string) (string, domain.OAuth2Token, error) {
	token, err := c.oauth2.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		return "", domain.OAuth2Token{}, err
	}
	userID, ok := token.Extra("user_id").(string)
	if !ok {
		return "", domain.OAuth2Token{}, fmt.Errorf("failed to get user_id")
	}

	return userID, domain.OAuth2Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       synchro.In[tz.AsiaTokyo](token.Expiry),
		CreatedAt:    synchro.Now[tz.AsiaTokyo](),
		UpdatedAt:    synchro.Now[tz.AsiaTokyo](),
	}, nil
}
