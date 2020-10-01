package authorizer

import (
	"github.com/graphql-go/graphql"
	"go.uber.org/zap"
)

type Authorizer interface {
	Authorize(resolver graphql.FieldResolveFn) graphql.FieldResolveFn
}

func NewAuthorizer(allowedRoles []string, logger *zap.Logger) Authorizer {
	return newAuthorizer(allowedRoles, logger)
}
