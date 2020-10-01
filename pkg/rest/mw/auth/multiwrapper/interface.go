package multiwrapper

import "github.com/graphql-go/graphql"

type MultiWrapper interface {
	Wrap(resolver graphql.FieldResolveFn) graphql.FieldResolveFn
}
