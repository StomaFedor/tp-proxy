package service

import (
	"context"
	"net/http"
	"net/url"
	"tp-proxy/pkg/models"
	"tp-proxy/pkg/repository"
)

type Request interface {
	SaveRequest(ctx context.Context,
		scheme,
		method,
		host,
		path string,
		headers http.Header,
		cookies []*http.Cookie,
		getParams url.Values,
		postParams url.Values) error
	GetAll(ctx context.Context) ([]models.Request, error)
	GetById(ctx context.Context, id int) (models.Request, error)
	RepeatRequest(request models.Request) (*http.Response, error)
	CheckSqlInjection(request models.Request) ([]string, error)
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
		Request:  NewRequestService(repo.Request),
		Responce: NewResponceService(repo.Responce),
	}
}
