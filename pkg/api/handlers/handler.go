package handler

import (
	"net/http"
	"tp-proxy/pkg/service"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "tp-proxy/docs"
)

type Handler struct {
	services *service.Service
}

func NewHandler(
	services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRoutes() http.Handler {
	r := mux.NewRouter()
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8000/swagger/doc.json"),
	))

	api := r.PathPrefix("/api").Subrouter()
	apiRouter := api.PathPrefix("/v1").Subrouter()

	apiRouter.HandleFunc("/requests", h.requests).Methods("GET")
	apiRouter.HandleFunc("/request/{id}", h.requestById).Methods("GET")
	// apiRouter.HandleFunc("/user", h.user).Methods("GET")
	// apiRouter.HandleFunc("/user", h.updateUser).Methods("POST", "OPTIONS")
	// apiRouter.HandleFunc("/user/share", h.getUserShareCridentials).Methods("GET")

	// api.Use(
	// 	h.loggingMiddleware,
	// 	h.panicRecoveryMiddleware,
	// 	h.corsMiddleware,
	// )

	return r
}
