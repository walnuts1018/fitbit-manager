package usecase

import "context"

func (u Usecase) GetName(ctx context.Context) (string, error) {
	return u.oauth2Client.GetName(ctx)
}
