package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"tp-proxy/pkg/repository"
	"tp-proxy/pkg/models"
)

type RequestService struct {
	RepoRequest repository.Request
}

func NewRequestService(repoRequest repository.Request) *RequestService {
	return &RequestService{RepoRequest: repoRequest}
}

func (s *RequestService) SaveRequest(ctx context.Context,
	method,
	host,
	path string,
	headers http.Header,
	cookies []*http.Cookie,
	getParams url.Values,
	postParams url.Values) error {

	request := models.Request{Method: method,
		Url:  path,
		Host: host,
	}
	request.Headers = make(map[string][]string)
	for key, value := range headers {
		request.Headers[key] = value
	}
	request.Cookies = make(map[string]string)
	for _, cookie := range cookies {
		request.Cookies[cookie.Name] = cookie.Value
	}
	request.GetParams = make(map[string][]string)
	for key, value := range getParams {
		request.GetParams[key] = value
	}

	request.PostParams = make(map[string][]string)
	for key, value := range postParams {
		request.PostParams[key] = value
	}

	reqJson, err := json.Marshal(request)
	if err != nil {
		return err
	}
	if _, err := s.RepoRequest.Save(ctx, reqJson); err != nil {
		return err
	}

	return nil
}

func (s *RequestService) GetAll(ctx context.Context) ([]models.Request, error) {
	return s.RepoRequest.GetAll(ctx)
}

func (s *RequestService) GetById(ctx context.Context, id int) (models.Request, error) {
	return s.RepoRequest.GetById(ctx, id)
}