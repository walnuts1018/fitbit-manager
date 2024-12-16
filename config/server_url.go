package config

import "net/url"

type ServerURL string

func ParseServerURL(v string) (ServerURL, error) {
	parsed, err := url.Parse(v)
	if err != nil {
		return "", err
	}
	return ServerURL(parsed.String()), nil
}
