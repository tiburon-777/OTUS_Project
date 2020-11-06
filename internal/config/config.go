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
		Capacity    int
		StoragePath string
	}
	Query struct {
		Timeout int
	}
	Log struct {
		File       string
		Level      string
		MuteStdout bool
	}
}

func NewConfig(configFile string) (Config, error) {
	var config Config
	f, err := os.Open(configFile)
	if err != nil {
		return config, err
	}
	defer f.Close()
	s, err := ioutil.ReadAll(f)
	if err != nil {
		return config, err
	}
	_, err = toml.Decode(string(s), &config)
	return config, err
}

func (c *Config) SetDefault() {
	c.Server = struct {
		Address string
		Port    string
	}{Address: "localhost", Port: "8080"}
	c.Cache = struct {
		Capacity    int
		StoragePath string
	}{Capacity: 20, StoragePath: "../assets/cache"}
	c.Query = struct{ Timeout int }{Timeout: 15}
	c.Log = struct {
		File       string
		Level      string
		MuteStdout bool
	}{File: "previewer.log", Level: "INFO", MuteStdout: false}
}
