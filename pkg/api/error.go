package api

import (
	"errors"
	"fmt"
)

type ErrorCode int

const (
	ContextDone ErrorCode = iota
	ContextNil
	ResolverNameNotFound
)

type Error struct {
	Code ErrorCode
	Msg  string
}

func (g *Error) Error() string {
	return fmt.Sprintf("%s: %d", g.Msg, g.Code)
}

func (g *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}

	return g.Code == t.Code
}

func HasErrorCode(err error, code ErrorCode) bool {
	return errors.Is(err, &Error{
		Code: code,
		Msg:  "",
	})
}
