package authorizer

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/sdeoras/graphql/pkg/rest/mw/auth"
	"go.uber.org/zap"
)

type authorizer struct {
	allowedRoles map[string]struct{}
	logger       *zap.Logger
}

func newAuthorizer(allowedRoles []string, logger *zap.Logger) *authorizer {
	g := &authorizer{
		allowedRoles: make(map[string]struct{}),
		logger:       logger,
	}

	for _, allowedRole := range allowedRoles {
		g.allowedRoles[allowedRole] = struct{}{}
	}

	return g
}

func (g *authorizer) Authorize(resolver graphql.FieldResolveFn) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		role, ok := p.Context.Value(auth.Role).(string)
		if !ok {
			return nil, fmt.Errorf("no valid role for the user")
		}
		if _, ok := g.allowedRoles[role]; !ok {
			return nil, fmt.Errorf("you don't have permission to access this field")
		}

		return resolver(p)
	}
}
