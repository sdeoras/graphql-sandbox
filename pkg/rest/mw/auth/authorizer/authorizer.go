package authorizer

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/sdeoras/graphql/pkg/rest/mw/auth"
	"go.uber.org/zap"
)

var rbac = map[string]map[string]map[string]struct{}{
	auth.RoleAdmin: {
		"*": auth.Permissions(
			auth.PermRead,
			auth.PermWrite,
			auth.PermUpdate,
			auth.PermDelete,
			auth.PermList,
		),
	},
	auth.RoleEditor: {
		"*": auth.Permissions(
			auth.PermRead,
			auth.PermWrite,
			auth.PermUpdate,
			auth.PermDelete,
			auth.PermList,
		),
	},
	auth.RoleViewer: {
		"*": auth.Permissions(
			auth.PermRead,
			auth.PermList,
		),
	},
}

type authorizer struct {
	permission string
	logger     *zap.Logger
}

func newAuthorizer(config *Config) *authorizer {
	g := &authorizer{
		permission: config.Permission,
		logger:     config.Logger,
	}

	return g
}

func (g *authorizer) Authorize(resolver graphql.FieldResolveFn) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		role, ok := p.Context.Value(auth.Role).(string)
		if !ok {
			return nil, fmt.Errorf("no valid role for the user")
		}

		paramPath := getParamPath(p.Info.Path.AsArray())

		rbacForRole, ok := rbac[role]
		if !ok {
			return nil, fmt.Errorf("your role does not have permissions")
		}

		for pathRegExp, permissions := range rbacForRole {
			if matched, err := filepath.Match(pathRegExp, paramPath); err == nil && matched {
				if _, ok := permissions[g.permission]; ok {
					return resolver(p)
				}
			}
		}

		return nil, fmt.Errorf("you don't have permission to access this field")
	}
}

func getParamPath(args []interface{}) string {
	var out []string
	for _, arg := range args {
		out = append(out, fmt.Sprintf("%s", arg))
	}
	return strings.Join(out, ".")
}
