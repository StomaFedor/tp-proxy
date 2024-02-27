package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

func (s *RequestService) RepeatRequest(request models.Request) (*http.Response, error) {
	req, err := makeHttpRequest(request)
	if err != nil {
		return nil, err
	}

	return http.DefaultTransport.RoundTrip(req)
}

func (s *RequestService) CheckSqlInjection(request models.Request) ([]string, error) {
	resp, err := s.RepeatRequest(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	status := resp.StatusCode
	contentLen := resp.ContentLength
	var result []string
	slice, err := checkSlice(request, status, contentLen, "get")
	if err != nil {
		return nil, err
	}
	result = append(result, slice...)
	slice, err = checkSlice(request, status, contentLen, "post")
	if err != nil {
		return nil, err
	}
	result = append(result, slice...)
	slice, err = checkSlice(request, status, contentLen, "headers")
	if err != nil {
		return nil, err
	}
	result = append(result, slice...)
	slice, err = checkCookie(request, status, contentLen)
	if err != nil {
		return nil, err
	}
	result = append(result, slice...)

	return result, nil
}

func checkSlice(request models.Request, status int, contentLen int64, checkType string) ([]string, error) {
	var result []string
	var collection map[string][]string
	switch checkType {
	case "get":
		collection = request.GetParams
	case "post":
		collection = request.PostParams
	case "headers":
		collection = request.Headers
	default:
		return nil, fmt.Errorf("invalid check type")
	}
	if len(collection) == 0 {
		return nil, nil
	}
	for name, item := range collection {
		for i, value := range item {
			newReq := request
			switch checkType {
			case "get":
				newReq.GetParams[name][i] = value + `'`
			case "post":
				newReq.PostParams[name][i] = value + `'`
			case "headers":
				newReq.Headers[name][i] = value + `'`
			}
			req, err := makeHttpRequest(newReq)
			if err != nil {
				return nil, err
			}
			resp, err := http.DefaultTransport.RoundTrip(req)
			if err != nil {
				return nil, err
			}
			if status != resp.StatusCode || contentLen != resp.ContentLength {
				result = append(result, fmt.Sprintf("%s %s: '%s' уязвим для SQL инъекций (одинарная кавычка)", checkType, name, value))
			}
			newReq = request
			switch checkType {
			case "get":
				newReq.GetParams[name][i] = value + `"`
			case "post":
				newReq.PostParams[name][i] = value + `"`
			case "headers":
				newReq.Headers[name][i] = value + "\""
			}
			req, err = makeHttpRequest(request)
			if err != nil {
				return nil, err
			}
			resp, err = http.DefaultTransport.RoundTrip(req)
			if err != nil {
				return nil, err
			}
			if status != resp.StatusCode || contentLen != resp.ContentLength {
				result = append(result, fmt.Sprintf("%s %s: '%s' уязвим для SQL инъекций (двойная кавычка)", checkType, name, value))
			}
		}
	}

	return result, nil
}

func checkCookie(request models.Request, status int, contentLen int64) ([]string, error) {
	var result []string
	for name, value := range request.Cookies {
		newReq := request
		newReq.Cookies[name] = value + `'`
		req, err := makeHttpRequest(request)
		if err != nil {
			return nil, err
		}
		resp, err := http.DefaultTransport.RoundTrip(req)
		if err != nil {
			return nil, err
		}
		if status != resp.StatusCode || contentLen != resp.ContentLength {
			result = append(result, fmt.Sprintf("cookie %s: '%s' уязвим для SQL инъекций (одинарная кавычка)", name, value))
		}
		newReq = request
		newReq.Cookies[name] = value + `"`
		req, err = makeHttpRequest(request)
		if err != nil {
			return nil, err
		}
		resp, err = http.DefaultTransport.RoundTrip(req)
		if err != nil {
			return nil, err
		}
		if status != resp.StatusCode || contentLen != resp.ContentLength {
			result = append(result, fmt.Sprintf("cookie %s: '%s' уязвим для SQL инъекций (двойная кавычка)", name, value))
		}
	}

	return result, nil
}

func makeHttpRequest(request models.Request) (*http.Request, error) {
	b, err := json.Marshal(request.PostParams)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(request.Method, request.Url, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	req.Host = request.Host
	req.URL.Scheme = "http"
	req.URL.Host = request.Host

	query := req.URL.Query()
	for key, value := range request.GetParams {
		for _, item := range value {
			query.Add(key, item)
		}
	}
	req.URL.RawQuery = query.Encode()

	for k, v := range request.Headers {
		for _, item := range v {
			req.Header.Set(k, fmt.Sprint(item))
		}
	}
	for k, v := range request.Cookies {
		req.AddCookie(&http.Cookie{Name: k, Value: fmt.Sprint(v)})
	}

	return req, nil
}
