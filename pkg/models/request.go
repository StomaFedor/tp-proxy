package models

type Request struct {
	Id         int
	Method     string
	Url        string
	Host       string
	Headers    map[string]any
	Cookies    map[string]any
	GetParams  map[string]any
	PostParams map[string]any
}
