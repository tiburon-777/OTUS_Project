package application

import (
	oslog "log"
	"net"
	"net/http"
	"time"

	"github.com/tiburon-777/OTUS_Project/previewer/config"
	"github.com/tiburon-777/OTUS_Project/previewer/logger"
)

type App struct {
	*http.Server
	Log logger.Interface
}

func New(conf config.Config) *App {
	loger, err := logger.New(conf.Log)
	if err != nil {
		oslog.Fatal("не удалось прикрутить логгер")
	}
	return &App{Server: &http.Server{Addr: net.JoinHostPort(conf.Server.Address, conf.Server.Port), Handler: LoggingMiddleware(http.HandlerFunc(Handler), loger)}, Log: loger}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	_, _ = w.Write([]byte("Hello! I'm calendar app!"))
}

func LoggingMiddleware(next http.Handler, l logger.Interface) http.HandlerFunc {
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
			l.Infof("receive %s request from IP: %s on path: %s, duration: %s useragent: %s ", r.Method, r.RemoteAddr, path, latency, useragent)
		}()
		next.ServeHTTP(w, r)
	})
}

func (s *App) Start() error {
	if err := s.ListenAndServe(); err != nil {
		return err
	}
	s.Log.Infof("Server starting")
	return nil
}

func (s *App) Stop() error {
	if err := s.Close(); err != nil {
		return err
	}
	s.Log.Infof("Server stoped")
	return nil
}
