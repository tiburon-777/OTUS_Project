package main

import (
	"flag"
	"github.com/tiburon-777/OTUS_Project/previewer/application"
	"github.com/tiburon-777/OTUS_Project/previewer/config"
	oslog "log"
	"os"
	"os/signal"
)

var ConfigFile = flag.String("config", "/etc/previewer.conf", "Path to configuration file")

func main() {
	flag.Parse()
	conf, err := config.NewConfig(*ConfigFile)
	if err != nil {
		oslog.Fatal("не удалось открыть файл конфигурации:", err.Error())
	}

	app := application.New(conf)

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals)

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