package domain

import "context"

type DataStore interface {
	RecordHeart(ctx context.Context, rates []HeartData) error
	GetLastHeartData(ctx context.Context) (HeartData, error)
	Close()
}
