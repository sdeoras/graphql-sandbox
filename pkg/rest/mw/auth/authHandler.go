package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/graphql-go/handler"
	"github.com/sdeoras/graphql/pkg/jwt"
	"go.uber.org/zap"
)

var Roles = map[string]string{
	GroupGoogle: RoleAdmin,
	GroupApple:  RoleViewer,
}

type authHandler struct {
	decoder jwt.Decoder
	handler *handler.Handler
	logger  *zap.Logger
}

type Config struct {
	PublicKey string
	Handler   *handler.Handler
	Logger    *zap.Logger
}

func NewHandler(config *Config) (http.Handler, error) {
	decoder, err := jwt.NewDecoder(config.PublicKey)
	if err != nil {
		return nil, err
	}
	return &authHandler{
		decoder: decoder,
		handler: config.Handler,
		logger:  config.Logger,
	}, nil
}

func (s *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := r.Header.Get(Authorization)
	s.logger.Debug("token", zap.String(Authorization, token))

	if token, err := getJwt(token); err == nil {
		s.logger.Debug("jwt", zap.String(XJwtToken, token))
		ctx = context.WithValue(ctx, XJwtToken, token)

		claims, err := s.decoder.Decode(token)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		email, ok := claims["email"].(string)
		if !ok {
			http.Error(w, "invalid email", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(email, "@")
		if len(parts) != 2 {
			http.Error(w, "invalid email", http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, XAuthenticated, true)
		ctx = context.WithValue(ctx, XUser, parts[0])
		ctx = context.WithValue(ctx, XGroups, []string{parts[1]})
		ctx = context.WithValue(ctx, Role, Roles[parts[1]])
		*r = *r.WithContext(ctx)
	} else {
		s.logger.Error("failed to get JWT token", zap.String("error", err.Error()))
		ctx = context.WithValue(ctx, XAuthenticated, false)
		*r = *r.WithContext(ctx)
	}

	s.handler.ServeHTTP(w, r)
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
