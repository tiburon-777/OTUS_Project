package application

import (
	"context"
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
		ctx, cansel := context.WithCancel(context.Background())
		defer cansel()
		q, err := buildQuery(r.URL)
		if err != nil {
			wErr := fmt.Errorf("can't parse query:\n %w", err)
			log.Warnf(wErr.Error())
			http.Error(w, wErr.Error(), http.StatusBadRequest)
			return
		}
		b, ok1, err := c.Get(cache.Key(q.id()))
		if err != nil {
			wErr := fmt.Errorf("can't get pic from cache:\n %w", err)
			log.Errorf(wErr.Error())
			http.Error(w, wErr.Error(), http.StatusInternalServerError)
			return
		}
		pic, ok2 := b.([]byte)
		if ok1 && ok2 {
			log.Infof("getting pic from cache")
			w.Header().Add("X-From-Appcache", "true")
			_, _ = w.Write(pic)
			return
		}
		pic, res, err := q.fromOrigin(ctx, r.Header, time.Duration(conf.Query.Timeout)*time.Second)
		if err != nil {
			wErr := fmt.Errorf("can't get pic from origin:\n %w", err)
			log.Warnf(wErr.Error())
			http.Error(w, wErr.Error(), http.StatusBadGateway)
			return
		}
		if res.StatusCode != 200 {
			log.Infof("Pic not found in origin or have problem with upstream. Response status:", res.Status)
			http.Error(w, "Pic not found in origin or have problem with upstream", res.StatusCode)
			return
		}
		pic, err = converter.SelectType(q.Width, q.Height, pic)
		if err != nil {
			wErr := fmt.Errorf("can't convert pic:\n %w", err)
			log.Errorf(wErr.Error())
			http.Error(w, wErr.Error(), http.StatusInternalServerError)
			return
		}
		_, err = c.Set(cache.Key(q.id()), pic)
		if err != nil {
			wErr := fmt.Errorf("can't add pic to cache:\n %w", err)
			log.Errorf(wErr.Error())
			http.Error(w, wErr.Error(), http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(pic)
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