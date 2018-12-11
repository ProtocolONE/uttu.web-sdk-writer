package types

type MsgNorm struct {
	ClientId    int64
	FingerPrint string

	//page
	PageUrl      string
	PageRef      string
	PageLanguage string
	PageCharset  string

	//body
	BodyIsIframe               bool
	BodyFramesCount            int64
	BodyIsSameOriginTopFrame   bool
	BodyIsInstantArticle       bool
	BodyNotificationPermission int64

	//screen
	ScreenMetricsBodyWidth        int64
	ScreenMetricsBodyHeight       int64
	ScreenMetricsDevicePixelRatio int64
	ScreenMetricsWidth            int64
	ScreenMetricsHeight           int64
	ScreenMetricsAvailableWight   int64
	ScreenMetricsAvailableHeight  int64
	ScreenMetricsOrientAngle      int64
	ScreenMetricsOrientType       string

	//time
	TimeStamp      int64
	TimeZone       string
	TimeZoneOffset int64

	//performance
	PerformanceDomainLookup                 int64
	PerformanceConnect                      int64
	PerformanceRequest                      int64
	PerformanceResponse                     int64
	PerformanceFetchFromNav                 int64
	PerformanceRedirect                     int64
	PerformanceRedirectCount                int64
	PerformanceDomInteractive               int64
	PerformanceDomContentLoadedEvent        int64
	PerformanceDomCompleteFromNav           int64
	PerformanceLoadEventStartFromNav        int64
	PerformanceLoadEvent                    int64
	PerformanceDomContentLoadedEventFromNav int64

	FirstContentfulPaint int64

	//platform
	PlatformIsJava bool
	PlatformFlash  string
	PlatformSystem string

	IsCookie bool

	//net
	NetDownlink      float32
	NetDownlinkMax   float32
	NetType          string
	NetEffectiveType string
	NetRTT           int64

	//tracker
	TrackerHitTime    int64
	TrackerNum        int64
	TrackerInactivity int64
	TrackerEventNum   int64
	TrackerIsExpired  bool
	TrackerId         string
}
