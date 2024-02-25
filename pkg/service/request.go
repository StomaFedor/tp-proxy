package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"tp-proxy/pkg/models"
	"tp-proxy/pkg/repository"
)

type RequestService struct {
	RepoRequest repository.Request
}

func NewRequestService(repoRequest repository.Request) *RequestService {
	return &RequestService{RepoRequest: repoRequest}
}

func (s *RequestService) SaveRequest(ctx context.Context,
	method,
	path string,
	headers http.Header,
	cookies []*http.Cookie,
	getParams url.Values,
	postParams url.Values) error {

	request := models.Request{Method: method,
		Url: path,
	}
	request.Headers = make(map[string]any)
	for key, value := range headers {
		request.Headers[key] = value
	}
	request.Cookies = make(map[string]any)
	for _, cookie := range cookies {
		request.Cookies[cookie.Name] = cookie.Value
	}
	request.GetPapams = make(map[string]any)
	for key, value := range getParams {
		request.GetPapams[key] = value
	}
	request.PostPapams = make(map[string]any)
	for key, value := range postParams {
		request.PostPapams[key] = value
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
