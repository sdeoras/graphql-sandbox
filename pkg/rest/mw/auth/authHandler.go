package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/graphql-go/handler"
	"go.uber.org/zap"
)

type authHandler struct {
	handler *handler.Handler
	logger  *zap.Logger
}

type Config struct {
	Handler *handler.Handler
	Logger  *zap.Logger
}

func NewHandler(config *Config) http.Handler {
	return &authHandler{
		handler: config.Handler,
		logger:  config.Logger,
	}
}

func (s *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer s.handler.ServeHTTP(w, r)

	ctx := r.Context()
	token := r.Header.Get(Authorization)
	s.logger.Debug("token", zap.String(Authorization, token))

	if jwt, err := getJwt(token); err == nil {
		s.logger.Debug("jwt", zap.String(XJwtToken, jwt))
		ctx = context.WithValue(ctx, Authorization, token)
		ctx = context.WithValue(ctx, XJwtToken, jwt)
		ctx = context.WithValue(ctx, XAuthenticated, true)
		ctx = context.WithValue(ctx, XGroups, []string{GroupGoogle}) // todo: decode from JWT
		ctx = context.WithValue(ctx, Role, RoleAdmin)                // todo
		*r = *r.WithContext(ctx)
	} else {
		s.logger.Error("failed to get JWT token", zap.String("error", err.Error()))
		ctx = context.WithValue(ctx, XAuthenticated, false)
		*r = *r.WithContext(ctx)
	}
}

func getJwt(token string) (string, error) {
	parts := strings.Split(token, " ")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid token format")
	}

	if strings.ToLower(parts[0]) != strings.ToLower(Bearer) {
		return "", fmt.Errorf("not a Bearer token")
	}

	return parts[1], nil
}
