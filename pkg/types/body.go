package types

type Body struct {
	BodyIframe                 int64  `json:"ifr"` //bool
	BodyFramesCount            int64  `json:"ifc"`
	BodySameOriginTopFrame     int64  `json:"sti"` //bool
	BodyInstantArticle         string `json:"iia"` //bool
	BodyNotificationPermission int64  `json:"ntf"`
}
