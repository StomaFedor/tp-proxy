package service

import (
	"context"
	"encoding/json"
	"net/http"
	"tp-proxy/pkg/models"
	"tp-proxy/pkg/repository"
)

type ResponceService struct {
	RepoResponce repository.Responce
}

func NewResponceService(repoResponce repository.Responce) *ResponceService {
	return &ResponceService{RepoResponce: repoResponce}
}

func (s *ResponceService) SaveResponce(ctx context.Context,
	code int,
	message string,
	headers http.Header,
	body []byte,
) error {

	responce := models.Responce{Code: code,
		Message: message,
	}
	responce.Headers = make(map[string]any)
	for key, value := range headers {
		responce.Headers[key] = value
	}
	responce.Body = string(body)

	respJson, err := json.Marshal(responce)
	if err != nil {
		return err
	}
	if _, err := s.RepoResponce.Save(ctx, respJson); err != nil {
		return err
	}

	return nil
}
