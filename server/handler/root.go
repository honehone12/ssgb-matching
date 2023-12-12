package handler

import (
	"net/http"
	"ssgb-matching/server/context"

	"github.com/labstack/echo/v4"
)

type RootResponse struct {
	Name    string
	Version string
}

func Root(c echo.Context) error {
	ctx, err := context.FromEchoContext(c)
	if err != nil {
		return err
	}

	m := ctx.Metadata()
	return c.JSON(http.StatusOK, RootResponse{
		Name:    m.Name(),
		Version: m.Version(),
	})
}
