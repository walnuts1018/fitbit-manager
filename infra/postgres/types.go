package postgres

import (
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/walnuts1018/fitbit-manager/domain"
	"gorm.io/gorm"
)

type OAuth2Token struct {
	UserID       string `gorm:"primarykey"`
	AccessToken  string
	RefreshToken string
	Expiry       synchro.Time[tz.AsiaTokyo]
	CreatedAt    synchro.Time[tz.AsiaTokyo]
	UpdatedAt    synchro.Time[tz.AsiaTokyo]
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func fromEntity(dto domain.OAuth2Token) OAuth2Token {
	return OAuth2Token{
		UserID:       dto.UserID,
		AccessToken:  dto.AccessToken,
		RefreshToken: dto.RefreshToken,
		Expiry:       dto.Expiry,
		CreatedAt:    dto.CreatedAt,
		UpdatedAt:    dto.UpdatedAt,
	}
}

func (t OAuth2Token) toEntity() domain.OAuth2Token {
	return domain.OAuth2Token{
		UserID:       t.UserID,
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
		Expiry:       t.Expiry,
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
	}
}
