package types

type Page struct {
	PageUrl      string `json:"url"`
	PageRef      string `json:"ref"`
	PageLanguage string `json:"la"`
	PageCharset  string `json:"cs"`
}
