package config

import "errors"

var ErrInvalidSessionSecretLength = errors.New("session secret must be 16, 24, or 32 bytes")

type CookieSecret string

func ParseCookieSecret(v string) (CookieSecret, error) {
	if len(v) != 16 && len(v) != 24 && len(v) != 32 {
		return "", ErrInvalidSessionSecretLength
	}
	return CookieSecret(v), nil
}
