package postgres

import (
	"context"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/walnuts1018/fitbit-manager/domain"
)

func (p *PostgresClient) SaveOAuth2Token(ctx context.Context, userID string, token domain.OAuth2Token) error {
	result := p.DB(ctx).Save(fromEntity(userID, token))
	if result.Error != nil {
		return fmt.Errorf("failed to update: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) GetOAuth2Token(ctx context.Context, userID string) (domain.OAuth2Token, error) {
	var token OAuth2Token
	result := p.DB(ctx).First(&token, userID)
	if result.Error != nil {
		return domain.OAuth2Token{}, fmt.Errorf("failed to get token: %w", result.Error)
	}
	return token.toEntity(), nil
}
