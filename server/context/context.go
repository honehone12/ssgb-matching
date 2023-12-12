package context

import (
	"ssgb-matching/errs"

	"github.com/labstack/echo/v4"
)

type Context struct {
	echo.Context
	*Components
}

func NewContext(e echo.Context, c *Components) *Context {
	return &Context{
		Context:    e,
		Components: c,
	}
}

func FromEchoContext(e echo.Context) (*Context, error) {
	ctx, ok := e.(*Context)
	if !ok {
		return nil, errs.ErrorCastFail("context")
	}
	return ctx, nil
}
