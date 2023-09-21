package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/walnuts1018/fitbit-manager/domain"
	"github.com/walnuts1018/fitbit-manager/mock/mockDatastore"
	"github.com/walnuts1018/fitbit-manager/mock/mockFitbit.go"
	"github.com/walnuts1018/fitbit-manager/mock/mockToken"
)

func TestUsecase_SignInAndCallback(t *testing.T) {
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

	oc := mockFitbit.NewMockOauth2Client()
	ts := mockToken.NewMockTokenClient()
	ds := mockDatastore.NewMockDatastore()

	tests := []struct {
		name    string
		fields  fields
		want    string
		want1   string
		wantErr bool
	}{
		{
			name: "normal",
			fields: fields{
				oauth2Client: oc,
				tokenStore:   ts,
				dataStore:    ds,
			},
			want:    "",
			want1:   "mockRedirectURL",
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
			_, got1, err := u.SignIn()
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.SignIn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got1 != tt.want1 {
				t.Errorf("Usecase.SignIn() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestUsecase_Callback(t *testing.T) {
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
		ctx  context.Context
		code string
	}

	oc := mockFitbit.NewMockOauth2Client()
	ts := mockToken.NewMockTokenClient()
	ds := mockDatastore.NewMockDatastore()

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
				ctx:  context.Background(),
				code: "mockCode",
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
			if err := u.Callback(tt.args.ctx, tt.args.code); (err != nil) != tt.wantErr {
				t.Errorf("Usecase.Callback() error = %v, wantErr %v", err, tt.wantErr)
			}

			token, err := tt.fields.tokenStore.GetOAuth2Token()
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.Callback() error = %v, wantErr %v", err, tt.wantErr)
			}
			if token.AccessToken != "mockAccessToken" {
				t.Errorf("Usecase.Callback() got = %v, want %v", token.AccessToken, "mockAccessToken")
			}
		})
	}
}
