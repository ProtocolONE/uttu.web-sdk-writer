package types

type Tracker struct {
	TrackerHitTime    int64  `json:"hitTime"`
	TrackerNum        int64  `json:"num"`
	TrackerInactivity int64  `json:"inactivity"`
	TrackerEventNum   int64  `json:"eventNum"`
	TrackerIsExpired  bool   `json:"isExpired"`
	TrackerId         string `json:"id"`
}
