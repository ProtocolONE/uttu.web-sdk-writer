package types

type Screen struct {
	ScreenMetricsBodyWidth        int64  `json:"w"`
	ScreenMetricsBodyHeight       int64  `json:"h"`
	ScreenMetricsDevicePixelRatio int64  `json:"dpr"`
	ScreenMetricsWidth            int64  `json:"tw"`
	ScreenMetricsHeight           int64  `json:"th"`
	ScreenMetricsAvailableWight   int64  `json:"aw"`
	ScreenMetricsAvailableHeight  int64  `json:"ah"`
	ScreenMetricsOrientAngle      int64  `json:"soa"`
	ScreenMetricsOrientType       string `json:"sot"`
}
