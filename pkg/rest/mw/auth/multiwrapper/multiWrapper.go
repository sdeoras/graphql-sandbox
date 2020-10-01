package multiwrapper

import (
	"github.com/sdeoras/graphql/pkg/rest/mw/auth/authenticator"
	"github.com/sdeoras/graphql/pkg/rest/mw/auth/authorizer"
)

type multiWrapper struct {
	authenticators []authenticator.Authenticator
	authorizers    []authorizer.Authorizer
}
