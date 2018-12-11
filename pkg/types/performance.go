package types

type Performance struct {
	PerformanceDomainLookup                 int64 `json:"dlp"`
	PerformanceConnect                      int64 `json:"ce"`
	PerformanceRequest                      int64 `json:"rqs"`
	PerformanceResponse                     int64 `json:"rss"`
	PerformanceFetchFromNav                 int64 `json:"fs"`
	PerformanceRedirect                     int64 `json:"re"`
	PerformanceRedirectCount                int64 `json:"rc"`
	PerformanceDomInteractive               int64 `json:"di"`
	PerformanceDomContentLoadedEvent        int64 `json:"dce"`
	PerformanceDomCompleteFromNav           int64 `json:"dc"`
	PerformanceLoadEventStartFromNav        int64 `json:"les"`
	PerformanceLoadEvent                    int64 `json:"le"`
	PerformanceDomContentLoadedEventFromNav int64 `json:"dcle"`
}
