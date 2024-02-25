package models

type Request struct {
	Method     string
	Url        string
	Headers    map[string]any
	Cookies    map[string]any
	GetPapams  map[string]any
	PostPapams map[string]any
}
