package domain

import (
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
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

type OAuth2Token struct {
	AccessToken  string
	RefreshToken string
	Expiry       synchro.Time[tz.AsiaTokyo]
	CreatedAt    synchro.Time[tz.AsiaTokyo]
	UpdatedAt    synchro.Time[tz.AsiaTokyo]
}
