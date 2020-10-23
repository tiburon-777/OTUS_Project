package application

import (
	"github.com/tiburon-777/OTUS_Project/previewer/cache"
	"github.com/tiburon-777/OTUS_Project/previewer/logger"
	"log"
	"net/http"
	"time"
)

func handler(c cache.Cache) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q,err := BuildQuery(r.URL)
		if err!=nil {
			http.Error(w, "Can't parse query", http.StatusNotFound)
			return
		}
		b,ok := c.Get(cache.Key(q.id()))
		pic,ok := b.([]byte)
		if ok {
			log.Println("Взяли из кэша")
			writeResponce(w, nil,200,pic)
			return
		}
		pic,h,err := q.fromOrigin()
		if err != nil {
			http.Error(w, "Pic not found in origin", http.StatusNotFound)
			return
		}
		pic, err = q.resize(pic)
		if err != nil {
			http.Error(w, "Resizer kirdyk...", http.StatusInternalServerError)
			return
		}
		c.Set(cache.Key(q.id()),pic)
		writeResponce(w, h,200,pic)
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

func writeResponce(w http.ResponseWriter, h http.Header, code int, body []byte) {
	for name, values := range h {
		for _, value := range values {
			w.Header().Add(name,value)
		}
	}
	_, _ = w.Write(body)
}

