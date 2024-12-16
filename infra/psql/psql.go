package psql

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/walnuts1018/fitbit-manager/config"
	"github.com/walnuts1018/fitbit-manager/domain"
)

const (
	sslMode = "disable"
)

type client struct {
	db *sql.DB
}

func NewPSQLClient(cfg config.Config) (client, error) {
	db, err := sql.Open("postgres", cfg.PSQLDSN.String())
	if err != nil {
		return client{}, fmt.Errorf("failed to open db: %v", err)
	}

	return client{db: db}, nil
}

func (c client) Close() error {
	return c.db.Close()
}

func (c client) SaveOAuth2Token(token domain.OAuth2Token) error {
	_, err := c.db.Exec(`CREATE TABLE IF NOT EXISTS oauth2_config (
		access_token TEXT NOT NULL,
		refresh_token TEXT NOT NULL,
		expiry TIMESTAMPTZ NOT NULL,
		created_at TIMESTAMPTZ NULL,
		updated_at TIMESTAMPTZ NULL
	)`)
	if err != nil {
		return fmt.Errorf("failed to create oauth2_config table: %v", err)
	}
	_, err = c.db.Exec("DELETE FROM oauth2_config")
	if err != nil {
		return fmt.Errorf("failed to delete oauth2_config: %v", err)
	}
	_, err = c.db.Exec("INSERT INTO oauth2_config (access_token, refresh_token, expiry, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)", token.AccessToken, token.RefreshToken, token.Expiry, token.CreatedAt, token.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert oauth2_config: %v", err)
	}
	return nil
}

func (c client) GetOAuth2Token() (domain.OAuth2Token, error) {
	var token domain.OAuth2Token
	err := c.db.QueryRow("SELECT access_token, refresh_token, expiry, created_at, updated_at FROM oauth2_config ORDER BY created_at DESC LIMIT 1;").Scan(&token.AccessToken, &token.RefreshToken, &token.Expiry, &token.CreatedAt, &token.UpdatedAt)
	if err != nil {
		return domain.OAuth2Token{}, fmt.Errorf("failed to get oauth2_config: %v", err)
	}
	return token, nil
}

func (c client) UpdateOAuth2Token(token domain.OAuth2Token) error {
	_, err := c.db.Exec("UPDATE oauth2_config SET access_token = $1, refresh_token = $2, expiry = $3, updated_at = $4", token.AccessToken, token.RefreshToken, token.Expiry, token.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update oauth2_config: %v", err)
	}
	return nil
}
