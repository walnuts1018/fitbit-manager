package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/walnuts1018/fitbit-manager/domain"
)

func (u *Usecase) GetHeart(ctx context.Context, from, to synchro.Time[tz.AsiaTokyo], detail domain.HeartDetail) ([]domain.HeartData, error) {
	hdatas := make([]domain.HeartData, 0)
	ttmp := from
	for to.After(ttmp) {
		if ttmp.Format("2006-01-02") == to.Format("2006-01-02") {
			hs, err := u.oauth2Client.GetHeartIntraday(ctx, ttmp.Format("2006-01-02"), ttmp.Format("15:04"), to.Format("15:04"), detail)
			if err != nil {
				return nil, fmt.Errorf("failed to get heart data: from:%v, to:%v, error: %v", ttmp, to, err)
			}
			for i := range hs {
				t, err := time.Parse("15:04:05", hs[i].Time)
				if err != nil {
					return nil, fmt.Errorf("failed to parse time, %v", hs[i].Time)
				}
				dt := time.Date(ttmp.Year(), ttmp.Month(), ttmp.Day(), t.Hour(), t.Minute(), 0, 0, timeJST.JST)
				hs[i].Datatime = &dt
			}
			hdatas = append(hdatas, hs...)
		} else {
			hs, err := u.oauth2Client.GetHeartIntraday(ctx, ttmp.Format("2006-01-02"), ttmp.Format("15:04"), "23:59", detail)
			if err != nil {
				return nil, fmt.Errorf("failed to get heart data: from:%v, to:%v, error:%v", ttmp, time.Date(ttmp.Year(), ttmp.Month(), ttmp.Day(), 23, 59, 0, 0, timeJST.JST), err)
			}
			for i := range hs {
				t, err := time.Parse("15:04:05", hs[i].Time)
				if err != nil {
					return nil, fmt.Errorf("failed to parse time, %v: %v", hs[i].Time, err)
				}
				dt := time.Date(ttmp.Year(), ttmp.Month(), ttmp.Day(), t.Hour(), t.Minute(), 0, 0, timeJST.JST)
				hs[i].Datatime = &dt
			}
			hdatas = append(hdatas, hs...)
		}
		ttmp = time.Date(ttmp.Year(), ttmp.Month(), ttmp.Day()+1, 0, 0, 0, 0, timeJST.JST)
	}
	return hdatas, nil
}

func (u *Usecase) GetHeartNow(ctx context.Context) (int, synchro.Time[tz.AsiaTokyo], error) {
	now := timeJST.Now()
	if u.heartCache.UpdatedAt.Add(1 * time.Minute).After(now) {
		slog.Info("use cache")
		return u.heartCache.heart, u.heartCache.dataAt, nil
	}

	hourBefore := now.Add(-6 * time.Hour)
	datas, err := u.GetHeart(ctx, hourBefore, now, domain.HeartDetailOneSecond)
	if err != nil {
		return 0, synchro.Time[tz.AsiaTokyo]{}, fmt.Errorf("failed to get heart rate: %w", err)
	}

	if len(datas) == 0 {
		return 0, synchro.Time[tz.AsiaTokyo]{}, fmt.Errorf("heart data is empty")
	}

	data := datas[len(datas)-1]
	if data.Datatime == nil {
		return 0, synchro.Time[tz.AsiaTokyo]{}, fmt.Errorf("datatime is nil")
	}
	u.heartCache.heart = data.Value
	u.heartCache.dataAt = *data.Datatime
	u.heartCache.UpdatedAt = now
	return data.Value, *data.Datatime, nil
}

func (u *Usecase) RecordHeart(ctx context.Context) error {
	now := timeJST.Now()
	last, err := u.dataStore.GetLastHeartData(ctx)
	if err != nil {
		return fmt.Errorf("failed to get last heart data: %w", err)
	}
	if last.Datatime == nil {
		slog.Warn("last datatime is nil")
		from := now.Add(-24 * time.Hour)
		last.Datatime = &from
	}
	data, err := u.GetHeart(ctx, *last.Datatime, now, domain.HeartDetailOneMinute)
	if err != nil {
		return fmt.Errorf("failed to get heart data: %w", err)
	}
	err = u.dataStore.RecordHeart(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to record heart data: %w", err)
	}
	return nil
}
