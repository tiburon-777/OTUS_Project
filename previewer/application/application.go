package application

import (
	"github.com/tiburon-777/OTUS_Project/previewer/cache"
	"github.com/tiburon-777/OTUS_Project/previewer/config"
	"github.com/tiburon-777/OTUS_Project/previewer/logger"
	oslog "log"
	"net"
	"net/http"
)

type App struct {
	*http.Server
	Log logger.Interface
	Cache cache.Cache
}

func New(conf config.Config) *App {
	loger, err := logger.New(conf.Log)
	if err != nil {
		oslog.Fatal("не удалось прикрутить логгер: ", err.Error())
	}
	c := cache.NewCache(conf.Cache.Capasity)
	return &App{Server: &http.Server{Addr: net.JoinHostPort(conf.Server.Address, conf.Server.Port), Handler: LoggingMiddleware(http.HandlerFunc(Handler), loger)}, Log: loger, Cache: c}
}

func (s *App) Start() error {
	s.Log.Infof("Server starting")
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

