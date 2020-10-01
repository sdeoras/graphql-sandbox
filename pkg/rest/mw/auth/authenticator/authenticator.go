package authenticator

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/sdeoras/graphql/pkg/rest/mw/auth"
	"go.uber.org/zap"
)

type authenticator struct {
	logger *zap.Logger
}

func newAuthenticator(logger *zap.Logger) *authenticator {
	return &authenticator{logger: logger}
}

func (s *authenticator) Authenticate(resolver graphql.FieldResolveFn) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		authenticated, ok := p.Context.Value(auth.XAuthenticated).(bool)
		if !ok || !authenticated {
			return nil, fmt.Errorf("not authenticated, ok:%v, authenticated:%v", ok, authenticated)
		}
		return resolver(p)
	}
}
