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
	skipCheck bool
	handler   *handler.Handler
	logger    *zap.Logger
}

type Config struct {
	Handler   *handler.Handler
	SkipCheck bool
	Logger    *zap.Logger
}

func NewHandler(config *Config) http.Handler {
	return &authHandler{
		handler:   config.Handler,
		skipCheck: config.SkipCheck,
		logger:    config.Logger,
	}
}

func (s *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	defer s.handler.ServeHTTP(w, r)

	if s.skipCheck {
		ctx = context.WithValue(ctx, XAuthenticated, true)
		ctx = context.WithValue(ctx, XGroups, []string{GroupGoogle})
		ctx = context.WithValue(ctx, Role, RoleAdmin)
		*r = *r.WithContext(ctx)
		return
	}

	token := r.Header.Get(Authorization)

	jwt, err := getJwt(token)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	ctx = context.WithValue(ctx, Authorization, token)
	ctx = context.WithValue(ctx, XJwtToken, jwt)
	ctx = context.WithValue(ctx, XAuthenticated, true)

	*r = *r.WithContext(ctx)
}

func getJwt(token string) (string, error) {
	parts := strings.Split(token, " ")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid token format")
	}

	if strings.ToLower(parts[0]) != Bearer {
		return "", fmt.Errorf("not a Bearer token")
	}

	return parts[1], nil
}
