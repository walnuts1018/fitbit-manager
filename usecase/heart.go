package usecase

import (
	"context"
	"time"

	"github.com/walnuts1018/fitbit-manager/domain"
	"github.com/walnuts1018/fitbit-manager/infra/timeJST"
)

func (u Usecase) GetHeartNow(ctx context.Context) (int, time.Time, error) {
	return u.oauth2Client.GetHeartNow(ctx)
}
func (u Usecase) RecordHeart(ctx context.Context) error {
	now := timeJST.Now()
	u.oauth2Client.GetHeart(ctx, now.Add(-12*time.Hour), now.Add(-6*time.Hour), domain.HeartDetailOneMinute)
	return nil
}
