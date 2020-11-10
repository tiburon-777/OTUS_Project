package application

import (
	"fmt"
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

func New(conf config.Config) (*App, error) {
	loger, err := logger.New(conf.Log)
	if err != nil {
		return nil, fmt.Errorf("can't start logger:\n %w", err)
	}
	c, err := cache.NewCache(conf.Cache.Capacity, conf.Cache.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("can't start cache:\n %w", err)
	}
	return &App{Server: &http.Server{Addr: net.JoinHostPort(conf.Server.Address, conf.Server.Port)}, Log: loger, Cache: c, Conf: conf}, nil
}

func (s *App) Start() error {
	s.Log.Infof("Server starting")
	s.Handler = loggingMiddleware(handler(s.Cache, s.Conf, s.Log), s.Log)
	err := s.ListenAndServe()
	s.Log.Infof("Server stoped")
	return err
}

func (s *App) Stop() error {
	if err := s.Close(); err != nil {
		return err
	}
	return nil
}
