package authenticator

import (
	"github.com/graphql-go/graphql"
	"go.uber.org/zap"
)

type Authenticator interface {
	Authenticate(resolver graphql.FieldResolveFn) graphql.FieldResolveFn
}

type Config struct {
	AllowedUsers  []string
	AllowedGroups []string
	Logger        *zap.Logger
}

func NewAuthenticator(config *Config) Authenticator {
	return newAuthenticator(config)
}
