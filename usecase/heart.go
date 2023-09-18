package usecase

import (
	"context"
	"time"
)

func (u Usecase) GetHeartNow(ctx context.Context) (int, time.Time, error) {
	return u.oauth2Client.GetHeartNow(ctx)
}
