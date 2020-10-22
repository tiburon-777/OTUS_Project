package application

import (
	"github.com/tiburon-777/OTUS_Project/previewer/logger"
	"net/http"
	"time"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	q,err := buildQuery(r.URL)
	if err!=nil {
		writeResponce(w, r.Header,501,[]byte("Can't parse query"))
	}
	pic, h, err := getPic(q)
	if err!=nil {
		writeResponce(w, h,501,[]byte("Have problem with cache"))
	}
	writeResponce(w, r.Header,200,pic)
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

func writeResponce(w http.ResponseWriter, h http.Header, code int, body []byte) {
	for name, values := range h {
		for _, value := range values {
			w.Header().Add(name,value)
		}
	}
	_, _ = w.Write(body)
}

