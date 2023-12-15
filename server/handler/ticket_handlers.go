package handler

import (
	"net/http"
	"ssgb-matching/conns"
	"ssgb-matching/errs"
	"ssgb-matching/matching/ticket"
	"ssgb-matching/server/context"
	"ssgb-matching/server/errres"
	"ssgb-matching/server/form"
	"ssgb-matching/uuid"

	"github.com/labstack/echo/v4"
)

type TicketNewForm struct {
	Class int64 `form:"class" validate:"required"`
}

type TicketNewResponse struct {
	Id string
}

func TicketNew(c echo.Context) error {
	ctx, err := context.FromEchoContext(c)
	if err != nil {
		return errres.InternalError(err, c.Logger())
	}

	formData := TicketNewForm{}
	if err := form.ProcessFormData(c, &formData); err != nil {
		return errres.BadRequest(err, c.Logger())
	}

	t, ch := ticket.MakeTicket(formData.Class)
	err = ctx.Engine().Enqueue(t)
	if errs.IsErrorIndexOutOfRange(err) {
		return errres.BadRequest(err, c.Logger())
	} else if err != nil {
		return errres.InternalError(err, c.Logger())
	}

	ctx.ConnMap().Set(t.Id(), conns.MakeConn(ctx.Engine().ConnParams(), ch, c.Logger()))

	return c.JSON(http.StatusOK, TicketNewResponse{
		Id: t.Id(),
	})
}

type TicketListenParams struct {
	Id string `validate:"required,uuid4,min=36,max=36"`
}

func TicketListen(c echo.Context) error {
	p := TicketListenParams{
		Id: c.Param("id"),
	}
	if err := c.Validate(&p); err != nil {
		return errres.BadRequest(err, c.Logger())
	}
	if err := uuid.ZeroUuid(p.Id); err != nil {
		return errres.BadRequest(err, c.Logger())
	}

	ctx, err := context.FromEchoContext(c)
	if err != nil {
		return errres.InternalError(err, c.Logger())
	}

	connMap := ctx.ConnMap()
	conn, err := connMap.Get(p.Id)
	if errs.IsErrorCastFail(err) {
		return errres.InternalError(err, c.Logger())
	} else if err != nil {
		return errres.BadRequest(err, c.Logger())
	}

	// !!
	if conn.Established() {
		return errres.BadRequest(errs.ErrorDuplicatedConnection(), c.Logger())
	}

	ws, err := ctx.WsUpgrader().Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return errres.InternalError(err, c.Logger())
	}

	conn.SetWs(ws)
	connMap.Set(p.Id, conn)

	return errs.ErrorNotimplemented()
}
