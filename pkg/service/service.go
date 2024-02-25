package service

import (
	"context"
	"net/http"
	"net/url"
	"tp-proxy/pkg/repository"
)

type Request interface {
	SaveRequest(ctx context.Context,
		method,
		path string,
		headers http.Header,
		cookies []*http.Cookie,
		getParams url.Values,
		postParams url.Values) error
	
}

type Responce interface {
	SaveResponce(ctx context.Context,
		code int,
		message string,
		headers http.Header,
		body []byte) error
	
}

type Service struct {
	Request  Request
	Responce Responce
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Request: NewRequestService(repo.Request),
		Responce: NewResponceService(repo.Responce),
	}
}
