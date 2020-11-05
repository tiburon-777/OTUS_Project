package application

import (
	"fmt"
	"net/http"
	"time"

	"github.com/tiburon-777/OTUS_Project/internal/cache"
	"github.com/tiburon-777/OTUS_Project/internal/config"
	"github.com/tiburon-777/OTUS_Project/internal/converter"
	"github.com/tiburon-777/OTUS_Project/internal/logger"
)

func handler(c cache.Cache, conf config.Config, log logger.Interface) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q, err := buildQuery(r.URL)
		if err != nil {
			wErr := fmt.Errorf("can't parse query: %w", err)
			log.Infof(wErr.Error())
			http.Error(w, wErr.Error(), http.StatusNotFound)
			return
		}
		b, ok1, err := c.Get(cache.Key(q.id()))
		if err != nil {
			wErr := fmt.Errorf("can't get pic from cache: %w", err)
			log.Infof(wErr.Error())
			http.Error(w, wErr.Error(), http.StatusInternalServerError)
			return
		}
		pic, ok2 := b.([]byte)
		if ok1 && ok2 {
			log.Infof("getting pic from cache")
			writeResponse(w, nil, pic)
			return
		}
		pic, res, err := q.fromOrigin(r.Header, time.Duration(conf.Query.Timeout)*time.Second)
		if err != nil {
			wErr := fmt.Errorf("can't get pic from origin: %w", err)
			log.Infof(wErr.Error())
			http.Error(w, wErr.Error(), http.StatusNotFound)
			return
		}
		if res.StatusCode != 200 {
			log.Infof("Pic not found in origin. Response status:", res.Status)
			http.Error(w, "Pic not found in origin", http.StatusNotFound)
			return
		}
		pic, err = converter.SelectType(q.Width, q.Height, pic)
		if err != nil {
			wErr := fmt.Errorf("can't convert pic: %w", err)
			log.Infof(wErr.Error())
			http.Error(w, wErr.Error(), http.StatusInternalServerError)
			return
		}
		_, err = c.Set(cache.Key(q.id()), pic)
		if err != nil {
			wErr := fmt.Errorf("can't add pic to cache: %w", err)
			log.Infof(wErr.Error())
			http.Error(w, wErr.Error(), http.StatusInternalServerError)
			return
		}
		writeResponse(w, nil, pic)
	})
}

func loggingMiddleware(next http.Handler, l logger.Interface) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			var path string
			if r.URL != nil {
				path = r.URL.Path
			}
			latency := time.Since(start)
			l.Infof("receive %s request from IP: %s on path: %s, duration: %s", r.Method, r.RemoteAddr, path, latency)
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
