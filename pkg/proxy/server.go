package proxy

import (
	"crypto/tls"
	"net/http"
	"time"
)

type ServerProxy struct {
	httpServer *http.Server

	protocol string
	cert string
	key string
}

func NewServerProxy(protoc, certif, kkey string) *ServerProxy {
	return &ServerProxy{
		protocol: protoc,
		cert: certif,
		key: kkey,
	}
}

func (s *ServerProxy) Serve(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	if s.protocol == "http" {
		return s.httpServer.ListenAndServe()
	} else {
		return s.httpServer.ListenAndServeTLS(s.cert, s.key)
	}
}
