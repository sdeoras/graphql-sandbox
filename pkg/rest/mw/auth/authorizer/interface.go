package authorizer

import (
	"github.com/graphql-go/graphql"
	"go.uber.org/zap"
)

type Authorizer interface {
	Authorize(resolver graphql.FieldResolveFn) graphql.FieldResolveFn
}

type Config struct {
	Permission string
	Logger     *zap.Logger
}

func NewAuthorizer(config *Config) Authorizer {
	return newAuthorizer(config)
}
