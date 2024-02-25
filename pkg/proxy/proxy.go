package proxy

import (
	"io"
	"log"
	"net"
	"net/http"
	"time"
	"tp-proxy/pkg/service"
)

type Proxy struct {
	services *service.Service
}

func NewProxy(services *service.Service) *Proxy {
	return &Proxy{services: services}
}

func (p *Proxy) InitHandlers() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodConnect {
			p.handleHTTPS(w, r)
		} else {
			p.handleHTTP(w, r)
		}
	})
}

func (p *Proxy) handleHTTP(w http.ResponseWriter, r *http.Request) {
	r.Header.Del("Proxy-Connection")
	r.RequestURI = ""

	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	err = r.ParseForm()
	if err != nil {
		log.Println("Failed to parse req form")
		return
	}

	err = p.services.Request.SaveRequest(r.Context(), r.Method, r.Host, r.URL.Path, r.Header, r.Cookies(), r.URL.Query(), r.Form)
	if err != nil {
		log.Println("Failed to save request")
		return
	}
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to parse resp data")
		return
	}
	err = p.services.Responce.SaveResponce(r.Context(), resp.StatusCode, resp.Status, resp.Header, data)
	if err != nil {
		log.Println("Failed to save responce")
		return
	}
}

func (p *Proxy) handleHTTPS(w http.ResponseWriter, r *http.Request) {
	dest_conn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	client_conn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	err = p.services.Request.SaveRequest(r.Context(), r.Method, r.Host, r.URL.Path, r.Header, r.Cookies(), r.URL.Query(), r.Form)
	if err != nil {
		log.Println("Failed to save request")
		return
	}
	go p.transfer(dest_conn, client_conn)
	go p.transfer(client_conn, dest_conn)
}

func (p *Proxy) transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

func copyHeader(dst, src http.Header) {
	for key, values := range src {
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}
