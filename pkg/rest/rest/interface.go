package rest

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

const (
	hello = "hello"
)

type Provider interface {
	Provide(name string) (http.Handler, error)
}

type Config struct {
	Logger *zap.Logger
}

func (g *Config) Validate() error {
	if g.Logger == nil {
		return fmt.Errorf("pl. provide a valid zap logger")
	}

	return nil
}

func NewProvider(config *Config) (Provider, error) {
	return newProvider(config)
}
