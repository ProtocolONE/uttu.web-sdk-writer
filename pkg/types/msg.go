package types

type Msg struct {
	Cookie               string `json:"c"` //bool
	FingerPrint          string `json:"fpr"`
	FirstContentfulPaint int64  `json:"fp"`

	Body        Body        `json:"bd"`
	Net         Net         `json:"nt"`
	Page        Page        `json:"page"`
	Performance Performance `json:"pm"`
	Platform    Platform    `json:"pl"`
	Screen      Screen      `json:"sc"`
	Time        Time        `json:"tm"`
	Tracker     []Tracker   `json:"tr"`
}
