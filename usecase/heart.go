package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/walnuts1018/fitbit-manager/domain"
	"github.com/walnuts1018/fitbit-manager/infra/timeJST"
)

func (u *Usecase) GetHeart(ctx context.Context, from, to time.Time, detail domain.HeartDetail) ([]domain.HeartData, error) {
	hdatas := make([]domain.HeartData, 0)
	ttmp := from
	for to.After(ttmp) {
		if ttmp.Format("2006-01-02") == to.Format("2006-01-02") {
			hs, err := u.oauth2Client.GetHeartIntraday(ttmp.Format("2006-01-02"), ttmp.Format("15:04"), to.Format("15:04"), detail)
			if err != nil {
				return nil, fmt.Errorf("failed to get heart data: from:%v, to:%v", ttmp, to)
			}
			for i := range hs {
				t, err := time.Parse("15:04", hs[i].Time)
				if err != nil {
					return nil, fmt.Errorf("failed to parse time, %v", hs[i].Time)
				}
				hs[i].Datatime = time.Date(ttmp.Year(), ttmp.Month(), ttmp.Day(), t.Hour(), t.Minute(), 0, 0, timeJST.JTC)
			}
			hdatas = append(hdatas, hs...)
		} else {
			hs, err := u.oauth2Client.GetHeartIntraday(ttmp.Format("2006-01-02"), ttmp.Format("15:04"), "23:59", detail)
			if err != nil {
				return nil, fmt.Errorf("failed to get heart data: from:%v, to:%v", ttmp, time.Date(ttmp.Year(), ttmp.Month(), ttmp.Day(), 23, 59, 0, 0, timeJST.JTC))
			}
			for i := range hs {
				t, err := time.Parse("15:04", hs[i].Time)
				if err != nil {
					return nil, fmt.Errorf("failed to parse time, %v", hs[i].Time)
				}
				hs[i].Datatime = time.Date(ttmp.Year(), ttmp.Month(), ttmp.Day(), t.Hour(), t.Minute(), 0, 0, timeJST.JTC)
			}
			hdatas = append(hdatas, hs...)
		}
		ttmp = time.Date(ttmp.Year(), ttmp.Month(), ttmp.Day()+1, 0, 0, 0, 0, timeJST.JTC)
	}
	return hdatas, nil
}

func (u *Usecase) GetHeartNow(ctx context.Context) (int, time.Time, error) {
	now := timeJST.Now()
	if u.heartCache.UpdatedAt.Add(1 * time.Minute).After(now) {
		slog.Info("use cache")
		return u.heartCache.heart, u.heartCache.dataAt, nil
	}

	hourBefore := now.Add(-1 * time.Hour)
	datas, err := u.GetHeart(ctx, hourBefore, now, domain.HeartDetailOneSecond)
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("failed to get heart rate: %w", err)
	}

	data := datas[len(datas)-1]
	dtime, err := time.Parse("2006-01-02 15:04:05", now.Format("2006-01-02 ")+data.Time)
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("failed to parse time: %w", err)
	}

	u.heartCache.heart = data.Value
	u.heartCache.dataAt = dtime
	u.heartCache.UpdatedAt = now
	return data.Value, dtime, nil
}

func (u *Usecase) RecordHeart(ctx context.Context) error {
	now := timeJST.Now()
	data, err := u.GetHeart(ctx, now.Add(-12*time.Hour), now.Add(-6*time.Hour), domain.HeartDetailOneMinute)
	if err != nil {
		return fmt.Errorf("failed to get heart data: %w", err)
	}
	err = u.dataStore.RecordHeart(data)
	if err != nil {
		return fmt.Errorf("failed to record heart data: %w", err)
	}
	return nil
}
