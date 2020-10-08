package config

import (
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Server struct {
		Address string
		Port    string
	}
	Cache struct {
		Capasity int
	}
	Log struct {
		File       string
		Level      string
		MuteStdout bool
	}
}

func NewConfig(configFile string) (Config, error) {
	f, err := os.Open(configFile)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()
	s, err := ioutil.ReadAll(f)
	if err != nil {
		return Config{}, err
	}
	var config Config
	_, err = toml.Decode(string(s), &config)
	return config, err
}
