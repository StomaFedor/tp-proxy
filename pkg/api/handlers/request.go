package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"tp-proxy/pkg/models"

	"github.com/gorilla/mux"
)

// @Summary get all requests
// @Tags requests
// @ID requests
// @Accept  json
// @Produce  json
// @Success 200 {object} ClientResponseDto[string]
// @Failure 500 {object} ClientResponseDto[string]
// @Router /api/v1/requests [get]
func (h *Handler) requests(w http.ResponseWriter, r *http.Request) {
	requests, err := h.services.Request.GetAll(r.Context())

	if err != nil {
		log.Println(err)
		NewErrorClientResponseDto(r.Context(), w, http.StatusInternalServerError, "internal server error")
		return
	}

	NewSuccessClientResponseDto(r.Context(), w, requests)
}

// @Summary get request by id
// @Tags request
// @ID request
// @Accept  json
// @Produce  json
// @Param id path string true "request id"
// @Success 200 {object} ClientResponseDto[string]
// @Failure 400,404,500 {object} ClientResponseDto[string]
// @Router /api/v1/request/{id} [get]
func (h *Handler) requestById(w http.ResponseWriter, r *http.Request) {
	idStr, ok := mux.Vars(r)["id"]
	if !ok {
		NewErrorClientResponseDto(r.Context(), w, http.StatusBadRequest, "invalid params")
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		NewErrorClientResponseDto(r.Context(), w, http.StatusBadRequest, "invalid params")
		return
	}

	request, err := h.services.Request.GetById(r.Context(), id)

	if err != nil {
		NewErrorClientResponseDto(r.Context(), w, http.StatusInternalServerError, "internal server error")
		return
	}

	NewSuccessClientResponseDto(r.Context(), w, request)
}

// @Summary get repeat request by id
// @Tags repeat
// @ID repeat
// @Accept  json
// @Produce  json
// @Param id path string true "request id"
// @Success 200 {object} ClientResponseDto[string]
// @Failure 400,404,500 {object} ClientResponseDto[string]
// @Router /api/v1/repeat/{id} [get]
func (h *Handler) repeatById(w http.ResponseWriter, r *http.Request) {
	idStr, ok := mux.Vars(r)["id"]
	if !ok {
		NewErrorClientResponseDto(r.Context(), w, http.StatusBadRequest, "invalid params")
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		NewErrorClientResponseDto(r.Context(), w, http.StatusBadRequest, "invalid params")
		return
	}

	request, err := h.services.Request.GetById(r.Context(), id)

	if err != nil {
		NewErrorClientResponseDto(r.Context(), w, http.StatusInternalServerError, "internal server error")
		return
	}

	data, err := h.repeatRequest(request)
	if err != nil {
		NewErrorClientResponseDto(r.Context(), w, http.StatusInternalServerError, "failed to repeat request")
		return
	}

	NewSuccessClientResponseDto(r.Context(), w, string(data))
}

func (h *Handler) repeatRequest(request models.Request) ([]byte, error) {
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
	for k, v := range request.Headers {
		for _, item := range v {
			req.Header.Set(k, fmt.Sprint(item))
		}
	}
	for k, v := range request.Cookies {
		req.AddCookie(&http.Cookie{Name: k, Value: fmt.Sprint(v)})
	}

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	h.services.Responce.SaveResponce(req.Context(), resp.StatusCode, resp.Status, resp.Header, data)
	return data, nil
}
