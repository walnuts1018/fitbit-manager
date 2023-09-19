package domain

type DataStore interface {
	RecordHeart(rates []HeartData) error
}
