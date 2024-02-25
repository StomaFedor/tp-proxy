package handler

import (
	"log"
	"net/http"
	"strconv"

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
		log.Println(err)
		NewErrorClientResponseDto(r.Context(), w, http.StatusInternalServerError, "internal server error")
		return
	}

	NewSuccessClientResponseDto(r.Context(), w, request)
}
