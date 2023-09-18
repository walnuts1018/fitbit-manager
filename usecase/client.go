package usecase

import "context"

func (u Usecase) NewFitbitClient(ctx context.Context) error {
	return u.oauth2Client.NewFitbitClient(ctx, u.tokenStore)
}
