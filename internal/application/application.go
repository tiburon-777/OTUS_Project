package application

import (
	oslog "log"
	"net"
	"net/http"

	"github.com/tiburon-777/OTUS_Project/internal/cache"
	"github.com/tiburon-777/OTUS_Project/internal/config"
	"github.com/tiburon-777/OTUS_Project/internal/logger"
)

type App struct {
	*http.Server
	Log   logger.Interface
	Cache cache.Cache
	Conf  config.Config
}

func New(conf config.Config) *App {
	loger, err := logger.New(conf.Log)
	if err != nil {
		oslog.Fatal("не удалось прикрутить логгер: ", err.Error())
	}
	c := cache.NewCache(conf.Cache.Capacity)
	return &App{Server: &http.Server{Addr: net.JoinHostPort(conf.Server.Address, conf.Server.Port)}, Log: loger, Cache: c, Conf: conf}
}

func (s *App) Start() error {
	s.Log.Infof("Server starting")
	s.Handler = loggingMiddleware(handler(s.Cache, s.Conf, s.Log), s.Log)
	_ = s.ListenAndServe()
	s.Log.Infof("Server stoped")
	return nil
}

func (s *App) Stop() error {
	if err := s.Close(); err != nil {
		return err
	}
	return nil
}
