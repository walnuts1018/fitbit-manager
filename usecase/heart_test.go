package usecase

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/walnuts1018/fitbit-manager/domain"
	"github.com/walnuts1018/fitbit-manager/infra/timeJST"
	"github.com/walnuts1018/fitbit-manager/mock/mockDatastore"
	"github.com/walnuts1018/fitbit-manager/mock/mockFitbit.go"
	"github.com/walnuts1018/fitbit-manager/mock/mockToken"
)

func TestUsecase_GetHeart(t *testing.T) {
	type fields struct {
		oauth2Client domain.FitbitClient
		tokenStore   domain.TokenStore
		dataStore    domain.DataStore
		heartCache   struct {
			heart     int
			dataAt    time.Time
			UpdatedAt time.Time
		}
	}
	type args struct {
		ctx    context.Context
		from   time.Time
		to     time.Time
		detail domain.HeartDetail
	}

	oc := mockFitbit.NewMockOauth2Client()
	ts := mockToken.NewMockTokenClient()
	ds := mockDatastore.NewMockDatastore()

	date1 := time.Date(2003, 10, 18, 0, 0, 0, 0, timeJST.JST)
	date2 := time.Date(2003, 10, 18, 6, 0, 0, 0, timeJST.JST)
	date3 := time.Date(2003, 10, 18, 23, 59, 0, 0, timeJST.JST)
	date4 := time.Date(2003, 10, 19, 0, 0, 0, 0, timeJST.JST)
	date5 := time.Date(2003, 10, 19, 23, 59, 0, 0, timeJST.JST)
	date6 := time.Date(2003, 10, 20, 0, 0, 0, 0, timeJST.JST)
	date7 := time.Date(2003, 10, 20, 6, 0, 0, 0, timeJST.JST)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.HeartData
		wantErr bool
	}{
		{
			name: "normal",
			fields: fields{
				oauth2Client: oc,
				tokenStore:   ts,
				dataStore:    ds,
			},
			args: args{
				ctx:    context.Background(),
				from:   date1,
				to:     date2,
				detail: domain.HeartDetailOneSecond,
			},
			want: []domain.HeartData{
				{
					Datatime: &date1,
					Time:     "00:00:00",
					Value:    100,
				},
				{
					Datatime: &date2,
					Time:     "06:00:00",
					Value:    100,
				},
			},
			wantErr: false,
		},
		{
			name: "multiday",
			fields: fields{
				oauth2Client: oc,
				tokenStore:   ts,
				dataStore:    ds,
			},
			args: args{
				ctx:    context.Background(),
				from:   date1,
				to:     date7,
				detail: domain.HeartDetailOneSecond,
			},
			want: []domain.HeartData{
				{
					Datatime: &date1,
					Time:     "00:00:00",
					Value:    100,
				},
				{
					Datatime: &date3,
					Time:     "23:59:00",
					Value:    100,
				},
				{
					Datatime: &date4,
					Time:     "00:00:00",
					Value:    100,
				},
				{
					Datatime: &date5,
					Time:     "23:59:00",
					Value:    100,
				},
				{
					Datatime: &date6,
					Time:     "00:00:00",
					Value:    100,
				},
				{
					Datatime: &date7,
					Time:     "06:00:00",
					Value:    100,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usecase{
				oauth2Client: tt.fields.oauth2Client,
				tokenStore:   tt.fields.tokenStore,
				dataStore:    tt.fields.dataStore,
				heartCache:   tt.fields.heartCache,
			}
			got, err := u.GetHeart(tt.args.ctx, tt.args.from, tt.args.to, tt.args.detail)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.GetHeart() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Usecase.GetHeart() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_GetHeartNow(t *testing.T) {
	type fields struct {
		oauth2Client domain.FitbitClient
		tokenStore   domain.TokenStore
		dataStore    domain.DataStore
		heartCache   struct {
			heart     int
			dataAt    time.Time
			UpdatedAt time.Time
		}
	}
	type args struct {
		ctx context.Context
	}

	oc := mockFitbit.NewMockOauth2Client()
	ts := mockToken.NewMockTokenClient()
	ds := mockDatastore.NewMockDatastore()
	timeJST.SetMockMode()
	now := timeJST.Now()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		want1   time.Time
		wantErr bool
	}{
		{
			name: "normal",
			fields: fields{
				oauth2Client: oc,
				tokenStore:   ts,
				dataStore:    ds,
			},
			args: args{
				ctx: context.Background(),
			},
			want:    100,
			want1:   time.Date(now.Year(), now.Month(), now.Day()-1, 23, 59, 0, 0, timeJST.JST),
			wantErr: false,
		},
		{
			name: "cached",
			fields: fields{
				oauth2Client: oc,
				tokenStore:   ts,
				dataStore:    ds,
				heartCache: struct {
					heart     int
					dataAt    time.Time
					UpdatedAt time.Time
				}{
					heart:     100,
					dataAt:    time.Date(now.Year(), now.Month(), now.Day()-1, 23, 58, 30, 0, timeJST.JST),
					UpdatedAt: now,
				},
			},
			args: args{
				ctx: context.Background(),
			},
			want:    100,
			want1:   time.Date(now.Year(), now.Month(), now.Day()-1, 23, 58, 30, 0, timeJST.JST),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usecase{
				oauth2Client: tt.fields.oauth2Client,
				tokenStore:   tt.fields.tokenStore,
				dataStore:    tt.fields.dataStore,
				heartCache:   tt.fields.heartCache,
			}
			got, got1, err := u.GetHeartNow(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.GetHeartNow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Usecase.GetHeartNow() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Usecase.GetHeartNow() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestUsecase_RecordHeart(t *testing.T) {
	type fields struct {
		oauth2Client domain.FitbitClient
		tokenStore   domain.TokenStore
		dataStore    domain.DataStore
		heartCache   struct {
			heart     int
			dataAt    time.Time
			UpdatedAt time.Time
		}
	}
	type args struct {
		ctx context.Context
	}

	oc := mockFitbit.NewMockOauth2Client()
	ts := mockToken.NewMockTokenClient()
	ds := mockDatastore.NewMockDatastore()
	timeJST.SetMockMode()

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "normal",
			fields: fields{
				oauth2Client: oc,
				tokenStore:   ts,
				dataStore:    ds,
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usecase{
				oauth2Client: tt.fields.oauth2Client,
				tokenStore:   tt.fields.tokenStore,
				dataStore:    tt.fields.dataStore,
				heartCache:   tt.fields.heartCache,
			}
			if err := u.RecordHeart(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Usecase.RecordHeart() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		hd, err := tt.fields.dataStore.GetLastHeartData(tt.args.ctx)
		if (err != nil) != tt.wantErr {
			t.Errorf("Usecase.RecordHeart() error = %v, wantErr %v", err, tt.wantErr)
		}
		if hd.Datatime == nil {
			t.Errorf("Usecase.RecordHeart() hd.Datatime is nil")
		}
		if hd.Value != 100 {
			t.Errorf("Usecase.RecordHeart() hd.Value = %v, want %v", hd.Value, 100)
		}
		if *hd.Datatime != time.Date(timeJST.Now().Year(), timeJST.Now().Month(), timeJST.Now().Day()-1, 23, 59, 0, 0, timeJST.JST) {
			t.Errorf("Usecase.RecordHeart() hd.Datatime = %v, want %v", *hd.Datatime, time.Date(timeJST.Now().Year(), timeJST.Now().Month(), timeJST.Now().Day()-1, 18, 00, 0, 0, timeJST.JST))
		}
	}
}
