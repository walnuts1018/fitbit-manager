package usecase

import (
	"context"
	"time"
)

func (u Usecase) GetHeartNow(ctx context.Context) (int, time.Time, error) {
	now := time.Now()
	return u.oauth2Client.GetHeart(ctx, now.Add(-1*time.Hour), now)
}
