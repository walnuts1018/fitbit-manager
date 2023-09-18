package domain

type TokenStore interface {
	SaveOAuth2Token(token OAuth2Token) error
	GetOAuth2Token() (OAuth2Token, error)
	Close() error
}
