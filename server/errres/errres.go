package errres

import (
	"net/http"
	"ssgb-matching/logger"

	"github.com/labstack/echo/v4"
)

func BadRequest(err error, l logger.Logger) error {
	l.Warn(err)
	return echo.NewHTTPError(http.StatusBadRequest, "invalid input")
}

func InternalError(err error, l logger.Logger) error {
	l.Error(err)
	return echo.NewHTTPError(http.StatusInternalServerError, "unexpected error")
}

func NotInService(l logger.Logger) error {
	l.Error("not in service")
	return echo.NewHTTPError(http.StatusInternalServerError, "not in service")
}
