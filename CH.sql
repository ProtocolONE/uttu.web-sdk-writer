CREATE TABLE WebSDKKafka
(
  ClientId                                Int64,
  FingerPrint                             String,

  PageUrl                                 String,
  PageRef                                 String,
  PageLanguage                            String,
  PageCharset                             String,

  BodyIsIframe                            UInt8,
  BodyFramesCount                         Int64,
  BodyIsSameOriginTopFrame                UInt8,
  BodyIsInstantArticle                    UInt8,
  BodyNotificationPermission              Int64,

  ScreenMetricsBodyWidth                  Int64,
  ScreenMetricsBodyHeight                 Int64,
  ScreenMetricsDevicePixelRatio           Int64,
  ScreenMetricsWidth                      Int64,
  ScreenMetricsHeight                     Int64,
  ScreenMetricsAvailableWight             Int64,
  ScreenMetricsAvailableHeight            Int64,
  ScreenMetricsOrientAngle                Int64,
  ScreenMetricsOrientType                 String,

  TimeStamp                               Int64,
  TimeZone                                String,
  TimeZoneOffset                          Int64,

  PerformanceDomainLookup                 Int64,
  PerformanceConnect                      Int64,
  PerformanceRequest                      Int64,
  PerformanceResponse                     Int64,
  PerformanceFetchFromNav                 Int64,
  PerformanceRedirect                     Int64,
  PerformanceRedirectCount                Int64,
  PerformanceDomInteractive               Int64,
  PerformanceDomContentLoadedEvent        Int64,
  PerformanceDomCompleteFromNav           Int64,
  PerformanceLoadEventStartFromNav        Int64,
  PerformanceLoadEvent                    Int64,
  PerformanceDomContentLoadedEventFromNav Int64,

  FirstContentfulPaint                    Int64,

  PlatformIsJava                          UInt8,
  PlatformFlash                           String,
  PlatformSystem                          String,

  IsCookie                                UInt8,

  NetDownlink                             Float32,
  NetDownlinkMax                          Float32,
  NetType                                 String,
  NetEffectiveType                        String,
  NetRTT                                  Int64,

  TrackerHitTime                          Int64,
  TrackerNum                              Int64,
  TrackerInactivity                       Int64,
  TrackerEventNum                         Int64,
  TrackerIsExpired                        UInt8,
  TrackerId                               String

) ENGINE = Kafka SETTINGS
  kafka_broker_list = '192.168.0.12:9092',
  kafka_topic_list = 'WebSDKPrepared-1-00,WebSDKPrepared-1-01,WebSDKPrepared-1-02',
  kafka_group_name = 'CHCluster',
  kafka_row_delimiter = '\n',
  kafka_format = 'JSONEachRow',
  kafka_num_consumers = 1;

CREATE TABLE WebSDK
(
  EventDate                               Date,
  ClientId                                Int64,
  FingerPrint                             String,

  PageUrl                                 String,
  PageRef                                 String,
  PageLanguage                            String,
  PageCharset                             String,

  BodyIsIframe                            UInt8,
  BodyFramesCount                         Int64,
  BodyIsSameOriginTopFrame                UInt8,
  BodyIsInstantArticle                    UInt8,
  BodyNotificationPermission              Int64,

  ScreenMetricsBodyWidth                  Int64,
  ScreenMetricsBodyHeight                 Int64,
  ScreenMetricsDevicePixelRatio           Int64,
  ScreenMetricsWidth                      Int64,
  ScreenMetricsHeight                     Int64,
  ScreenMetricsAvailableWight             Int64,
  ScreenMetricsAvailableHeight            Int64,
  ScreenMetricsOrientAngle                Int64,
  ScreenMetricsOrientType                 String,

  TimeStamp                               Int64,
  TimeZone                                String,
  TimeZoneOffset                          Int64,

  PerformanceDomainLookup                 Int64,
  PerformanceConnect                      Int64,
  PerformanceRequest                      Int64,
  PerformanceResponse                     Int64,
  PerformanceFetchFromNav                 Int64,
  PerformanceRedirect                     Int64,
  PerformanceRedirectCount                Int64,
  PerformanceDomInteractive               Int64,
  PerformanceDomContentLoadedEvent        Int64,
  PerformanceDomCompleteFromNav           Int64,
  PerformanceLoadEventStartFromNav        Int64,
  PerformanceLoadEvent                    Int64,
  PerformanceDomContentLoadedEventFromNav Int64,

  FirstContentfulPaint                    Int64,

  PlatformIsJava                          UInt8,
  PlatformFlash                           String,
  PlatformSystem                          String,

  IsCookie                                UInt8,

  NetDownlink                             Float32,
  NetDownlinkMax                          Float32,
  NetType                                 String,
  NetEffectiveType                        String,
  NetRTT                                  Int64,

  TrackerHitTime                          Int64,
  TrackerNum                              Int64,
  TrackerInactivity                       Int64,
  TrackerEventNum                         Int64,
  TrackerIsExpired                        UInt8,
  TrackerId                               String
) ENGINE = MergeTree
  (
  )
    PARTITION BY toYYYYMM
      (
        EventDate
      )
    ORDER BY
      (
       ClientId,
       EventDate,
       FingerPrint
        )
    SETTINGS index_granularity = 8192;

CREATE
  MATERIALIZED
  VIEW
  WebSDKKafkaMV
  TO
    WebSDK
AS
SELECT toDate(round(divide(TimeStamp, 1000))) AS EventDate,*
FROM WebSDKKafka;