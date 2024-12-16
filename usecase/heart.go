package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/walnuts1018/fitbit-manager/domain"
)

func (u *Usecase) GetHeartNow(ctx context.Context, userID string) (domain.HeartData, error) {
	data, err := u.dataStore.GetLatestHeartData(ctx, userID)
	if err != nil {
		return domain.HeartData{}, fmt.Errorf("failed to get last heart data: %w", err)
	}
	return data, nil
}

func (u *Usecase) RecordHeart(ctx context.Context, userID string) error {
	var from synchro.Time[tz.AsiaTokyo]

	now := synchro.Now[tz.AsiaTokyo]()

	latest, err := u.dataStore.GetLatestHeartData(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			from = now.Add(-24 * time.Hour)
		} else {
			return fmt.Errorf("failed to get last heart data: %w", err)
		}
	} else {
		from = latest.Time
	}

	token, err := u.tokenStore.GetOAuth2Token(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get oauth2 token: %w", err)
	}

	timeRanges := domain.NewFitbitTimeRange(from, now)
	for _, r := range timeRanges {
		data, newToken, err := u.fitbitClient.GetHeartData(ctx, token, r, domain.HeartDetailOneMinute)
		if err != nil {
			return fmt.Errorf("failed to get heart data: %w", err)
		}
		token = newToken

		if err := u.dataStore.RecordHeart(ctx, userID, data); err != nil {
			return fmt.Errorf("failed to record heart data: %w", err)
		}
	}

	if err := u.tokenStore.SaveOAuth2Token(ctx, userID, token); err != nil {
		return fmt.Errorf("failed to save oauth2 token: %w", err)
	}

	return nil
}
