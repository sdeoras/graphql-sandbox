package resolvers

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"go.uber.org/zap"
)

type Resolver interface {
	Resolve(name string) (
		func(p graphql.ResolveParams) (interface{}, error),
		error,
	)
}

type Config struct {
	Logger *zap.Logger
}

func (g *Config) Validate() error {
	if g.Logger == nil {
		return fmt.Errorf("please provide a valid zap logger instance")
	}

	return nil
}

func NewResolver(config *Config) (Resolver, error) {
	return newResolver(config)
}
