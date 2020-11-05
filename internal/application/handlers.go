package application

import (
	"log"
	"net/http"
	"time"

	"github.com/tiburon-777/OTUS_Project/internal/cache"
	"github.com/tiburon-777/OTUS_Project/internal/config"
	"github.com/tiburon-777/OTUS_Project/internal/converter"
	"github.com/tiburon-777/OTUS_Project/internal/logger"
)

func handler(c cache.Cache, conf config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q, err := buildQuery(r.URL)
		if err != nil {
			http.Error(w, "Can't parse query", http.StatusNotFound)
			return
		}
		b, ok1 := c.Get(cache.Key(q.id()))
		pic, ok2 := b.([]byte)
		if ok1 && ok2 {
			log.Println("Взяли из кэша")
			writeResponse(w, nil, pic)
			return
		}
		pic, res, err := q.fromOrigin(time.Duration(conf.Query.Timeout) * time.Second)
		if err != nil || res.StatusCode != 200 {
			http.Error(w, "Pic not found in origin", http.StatusNotFound)
			return
		}
		pic, err = converter.SelectType(q.Width, q.Height, pic)
		if err != nil {
			http.Error(w, "Resizer kirdyk...", http.StatusInternalServerError)
			return
		}
		c.Set(cache.Key(q.id()), pic)
		writeResponse(w, nil, pic)
	})
}

func loggingMiddleware(next http.Handler, l logger.Interface) http.HandlerFunc {
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

func writeResponse(w http.ResponseWriter, h http.Header, body []byte) {
	for name, values := range h {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}
	_, _ = w.Write(body)
}
