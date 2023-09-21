package timeJST

import "time"

var (
	JST      *time.Location
	mockMode = false
	MockTime time.Time
)

func init() {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	JST = jst
}

func Now() time.Time {
	if mockMode {
		return time.Date(2003, 10, 18, 0, 0, 0, 0, JST)
	}
	return time.Now().In(JST)
}

func SetMockMode() {
	mockMode = true
}
