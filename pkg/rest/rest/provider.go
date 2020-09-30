package rest

import (
	"fmt"
	"net/http"

	"github.com/sdeoras/graphql/pkg/rest/api"
	"go.uber.org/zap"
)

type provider struct {
	logger    *zap.Logger
	httpError func(w http.ResponseWriter, msg string, code int)
}

func newProvider(config *Config) (*provider, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &provider{
		logger: config.Logger,
		httpError: func(w http.ResponseWriter, msg string, code int) {
			http.Error(w, msg, code)
		},
	}, nil
}

func (g *provider) Provide(name string) (http.Handler, error) {
	switch name {
	case hello:
		return http.HandlerFunc(g.hello()), nil
	default:
		return nil, fmt.Errorf("invalid name")
	}
}

func (g *provider) hello() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		request := &api.GraphQlRequest{}
		if err := request.ReadFrom(r); err != nil {
			msg := fmt.Errorf("error unmarshaling request:%w", err)
			g.logger.Error(msg.Error())
			g.httpError(w, "invalid payload", http.StatusBadRequest)
			return
		}

	}
}
