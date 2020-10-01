package authenticator

import (
	"github.com/graphql-go/graphql"
	"go.uber.org/zap"
)

type Authenticator interface {
	Authenticate(resolver graphql.FieldResolveFn) graphql.FieldResolveFn
}

func NewAuthenticator(logger *zap.Logger) Authenticator {
	return newAuthenticator(logger)
}
