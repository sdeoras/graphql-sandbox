package handler

import (
	"context"
	"fmt"
	"github.com/graphql-go/handler"
	"github.com/sdeoras/graphql/pkg/rest/mw/auth"
	"net/http"
	"strings"
)

type authHandler struct {
	skipCheck bool
	h         *handler.Handler
}

func New(h *handler.Handler, skipCheck bool) http.Handler {
	return &authHandler{
		h:         h,
		skipCheck: skipCheck,
	}
}

func (s *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	defer s.h.ServeHTTP(w, r)

	if s.skipCheck {
		ctx = context.WithValue(ctx, auth.XAuthenticated, true)
		ctx = context.WithValue(ctx, auth.XGroups, []string{auth.GroupGoogle})
		ctx = context.WithValue(ctx, auth.Role, auth.RoleAdmin)
		*r = *r.WithContext(ctx)
		return
	}

	token := r.Header.Get(auth.Authorization)

	jwt, err := getJwt(token)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	ctx = context.WithValue(ctx, auth.Authorization, token)
	ctx = context.WithValue(ctx, auth.XJwtToken, jwt)
	ctx = context.WithValue(ctx, auth.XAuthenticated, true)

	*r = *r.WithContext(ctx)
}

func getJwt(token string) (string, error) {
	parts := strings.Split(token, " ")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid token format")
	}

	if strings.ToLower(parts[0]) != auth.Bearer {
		return "", fmt.Errorf("not a Bearer token")
	}

	return parts[1], nil
}
