package server

import (
	"ssgb-matching/server/context"
	"ssgb-matching/server/handler"
	"ssgb-matching/server/validator"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echolog "github.com/labstack/gommon/log"
)

type ServerParams struct {
	ListenAt string
	LogLevel echolog.Lvl
}

type Server struct {
	params     ServerParams
	echo       *echo.Echo
	components *context.Components
	errCh      chan error
}

func NewServer(
	params ServerParams,
	e *echo.Echo,
	c *context.Components,
) *Server {
	e.Validator = validator.NewValidator()
	e.Logger.SetLevel(params.LogLevel)
	return &Server{
		params:     params,
		echo:       e,
		components: c,
		errCh:      make(chan error),
	}
}

func (s *Server) convertContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := context.NewContext(c, s.components)
		return next(ctx)
	}
}

func (s *Server) start() {
	s.errCh <- s.echo.Start(s.params.ListenAt)
}

func (s *Server) Run() <-chan error {
	s.echo.Use(s.convertContext)
	s.echo.Use(middleware.Recover())
	s.echo.Use(middleware.Logger())

	s.echo.GET("/", handler.Root)
	s.echo.POST("/ticket/new", handler.TicketNew)
	s.echo.GET("/ticket/listen/;id", handler.TicketListen)

	go s.start()
	return s.errCh
}
