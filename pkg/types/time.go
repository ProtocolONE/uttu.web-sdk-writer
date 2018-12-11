package types

type Time struct {
	TimeStamp      int64  `json:"ts"`
	TimeZone       string `json:"tz"`
	TimeZoneOffset int64  `json:"tzo"`
}
