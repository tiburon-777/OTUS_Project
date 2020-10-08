package main
import (
	"flag"
	oslog "log"
	"github/tiburon-777/OTUS_Project/previewer/application"
)

var ConfigFile = flag.String("config", "/etc/previewer.conf", "Path to configuration file")

func main() {
	flag.Parse()
	conf, err := config.NewConfig(*ConfigFile)
	if err != nil {
		oslog.Fatal("не удалось открыть файл конфигурации:", err.Error())
	}

	app := application.New(&conf)
	app.Run()
}