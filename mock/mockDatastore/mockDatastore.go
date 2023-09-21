package mockDatastore

import (
	"context"

	"github.com/walnuts1018/fitbit-manager/domain"
)

type mockDatastore struct {
	heartDatas []domain.HeartData
}

func NewMockDatastore() domain.DataStore {
	return &mockDatastore{}
}

func (c *mockDatastore) RecordHeart(ctx context.Context, rates []domain.HeartData) error {
	c.heartDatas = append(c.heartDatas, rates...)
	return nil
}

func (c *mockDatastore) GetLastHeartData(ctx context.Context) (domain.HeartData, error) {
	if len(c.heartDatas) == 0 {
		return domain.HeartData{}, nil
	}
	return c.heartDatas[len(c.heartDatas)-1], nil
}

func (c *mockDatastore) Close() {}
