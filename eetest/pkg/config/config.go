package config

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
)

type Config struct {
	Host         string        `long:"host" env:"HOST" default:"0.0.0.0" description:"The host for the service"`
	Port         int           `long:"port" env:"PORT" default:"8080" description:"The port the collector service is listening on"`
	WriteTimeout int         `long:"write_timeout" env:"WRITE_TIMEOUT" default:"15" description:"The server WriteTimeout"`
	ReadTimeout  int `long:"read_timeout" env:"READ_TIMEOUT" default:"15" description:"The server ReadTimeout"`
	IdleTimeout  int `long:"idle_timeout" env:"IDLE_TIMEOUT" default:"15" description:"The server IdleTimeout"`
}

func NewConfig() (*Config, error) {
	conf := &Config{}
	parser := flags.NewParser(conf, flags.Default)
	if _, err := parser.Parse(); err != nil {
		parser.WriteHelp(os.Stderr)
		return nil, fmt.Errorf("Failed to parse config: %v", err)
	}
	return conf, nil
}
