package authenticator

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/sdeoras/graphql/pkg/rest/mw/auth"
	"go.uber.org/zap"
)

type authenticator struct {
	allowedUsers  map[string]struct{}
	allowedGroups map[string]struct{}
	logger        *zap.Logger
}

func newAuthenticator(config *Config) *authenticator {
	g := &authenticator{
		allowedUsers:  map[string]struct{}{},
		allowedGroups: map[string]struct{}{},
		logger:        config.Logger,
	}

	for _, allowedUser := range config.AllowedUsers {
		g.allowedUsers[allowedUser] = struct{}{}
	}

	for _, allowedGroup := range config.AllowedGroups {
		g.allowedGroups[allowedGroup] = struct{}{}
	}

	return g
}

func (s *authenticator) Authenticate(resolver graphql.FieldResolveFn) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		authenticated, ok := p.Context.Value(auth.XAuthenticated).(bool)
		if !ok || !authenticated {
			return nil, fmt.Errorf("not authenticated")
		}

		groupMatched, userMatched := false, false
		if groups, ok := p.Context.Value(auth.XGroups).([]string); ok {
			for _, group := range groups {
				if _, ok := s.allowedGroups[group]; ok {
					groupMatched = true
					break
				}
			}
		}

		if !groupMatched {
			if user, ok := p.Context.Value(auth.XUser).(string); ok {
				if _, ok := s.allowedUsers[user]; ok {
					userMatched = true
				}
			}
		}

		if groupMatched || userMatched {
			return resolver(p)
		}

		return nil, fmt.Errorf("not authenticated")
	}
}
