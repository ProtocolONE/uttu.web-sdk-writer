package types

type Net struct {
	NetDownlink      float32 `json:"dl"`
	NetDownlinkMax   float32 `json:"dm"`
	NetType          string  `json:"t"`
	NetEffectiveType string  `json:"et"`
	NetRTT           int64   `json:"rtt"`
}
