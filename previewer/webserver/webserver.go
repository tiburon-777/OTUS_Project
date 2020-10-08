package webserver

import (
	"net"
	"net/http"
	"time"
)

type Server struct {
	*http.Server
}

func NewServer(address string, port string) http.Server {
	return http.Server{Addr: net.JoinHostPort(address, port), Handler: LoggingMiddleware(Handler)}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	_, _ = w.Write([]byte("Hello! I'm calendar app!"))
}

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			var path, useragent string
			if r.URL != nil {
				path = r.URL.Path
			}
			if len(r.UserAgent()) > 0 {
				useragent = r.UserAgent()
			}
			latency := time.Since(start)
			a.Logger.Infof("receive %s request from IP: %s on path: %s, duration: %s useragent: %s ", r.Method, r.RemoteAddr, path, latency, useragent)
		}()
		next.ServeHTTP(w, r)
	})
}

func (s *Server) Start() error {
	if err := s.ListenAndServe(); err != nil {
		return err
	}
	s.app.Logger.Infof("Server starting")
	return nil
}

func (s *Server) Stop() error {
	if err := s.Close(); err != nil {
		return err
	}
	s.app.Logger.Infof("Server stoped")
	return nil
}
