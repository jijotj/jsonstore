package config

import (
	"fmt"
)

const (
	servePortConfKey     = "SERVE_PORT"
	serverTimeOutConfKey = "SERVER_TIMEOUT_MS"
)

type Config struct {
	ServePort       int
	ServerTimeoutMS int
}

func New() (*Config, error) {
	vars := Vars{}
	serverPort := vars.MandatoryInt(servePortConfKey)
	serverTimeoutMs := vars.OptionalInt(serverTimeOutConfKey, 10)

	if err := vars.Error(); err != nil {
		return nil, fmt.Errorf("config: environment variables: %s", err)
	}

	return &Config{
		ServePort:       serverPort,
		ServerTimeoutMS: serverTimeoutMs,
	}, nil
}
