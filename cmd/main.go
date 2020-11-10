package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tiburon-777/OTUS_Project/internal/application"
	"github.com/tiburon-777/OTUS_Project/internal/config"
)

var (
	ConfigFile = flag.String("config", "/etc/previewer.conf", "Path to configuration file")
	CleanCache = flag.Bool("clean", false, "Set true if you need clean cache before start app")
)

func main() {
	flag.Parse()
	conf, err := config.NewConfig(*ConfigFile)
	if err != nil {
		log.Println("Configuration file not found. Will use defaults.")
		conf.SetDefault()
	}
	app, err := application.New(conf)
	if err != nil {
		log.Fatal("can't start application:", err.Error())
	}
	if *CleanCache {
		err = app.Cache.Clear()
		if err != nil {
			log.Fatalf("can't clean cache:\n %s", err.Error())
		}
	}
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

		<-signals
		signal.Stop(signals)

		if err := app.Stop(); err != nil {
			app.Log.Errorf("failed to close application: " + err.Error())
		}
	}()

	if err := app.Start(); err != nil {
		app.Log.Errorf("failed to start application: " + err.Error())
		os.Exit(1)
	}
}
