package resolvers

import (
	"github.com/graphql-go/graphql"
	"github.com/sdeoras/graphql/pkg/api"
	"go.uber.org/zap"
)

const (
	Hello = "hello"
)

type resolver struct {
	logger *zap.Logger
}

func newResolver(config *Config) (*resolver, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &resolver{
		logger: config.Logger,
	}, nil
}

func (g *resolver) Resolve(name string) (func(p graphql.ResolveParams) (interface{}, error), error) {
	switch name {
	case Hello:
		return g.hello(), nil
	default:
		msg := "resolver name not found"
		g.logger.Error(msg)
		return nil, &api.Error{
			Code: api.ResolverNameNotFound,
			Msg:  msg,
		}
	}
}

func (g *resolver) hello() func(p graphql.ResolveParams) (interface{}, error) {
	return func(p graphql.ResolveParams) (interface{}, error) {

		if p.Context == nil {
			return nil, &api.Error{
				Code: api.ContextNil,
				Msg:  "input context must be provided",
			}
		}

		// always check context
		select {
		case <-p.Context.Done():
			return nil, &api.Error{
				Code: api.ContextDone,
				Msg:  "input context is done",
			}
		default:
		}

		return "world", nil
	}
}
