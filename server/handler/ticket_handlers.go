package handler

import (
	"net/http"
	"ssgb-matching/conns"
	"ssgb-matching/errs"
	"ssgb-matching/gsip"
	"ssgb-matching/matching/tickets"
	"ssgb-matching/server/context"
	"ssgb-matching/server/errres"
	"ssgb-matching/server/form"
	"ssgb-matching/uuid"
	"strconv"

	"github.com/labstack/echo/v4"
)

type TicketNewForm struct {
	Class    int64  `form:"class" validate:"required"`
	Backfill string `form:"backfill" validate:"required,boolean"`
}

type TicketNewResponse struct {
	Id string

	FoundBackfill bool
	BackfillGsip  gsip.GSIP
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

	useBackfill, err := strconv.ParseBool(formData.Backfill)
	if err != nil {
		return errres.BadRequest(err, c.Logger())
	}

	engine := ctx.Engine()

	if useBackfill {
		backfill, err := engine.FindBackfill(formData.Class)
		if !errs.IsErrorNotFound(err) {
			if errs.IsErrorIndexOutOfRange(err) {
				return errres.BadRequest(err, c.Logger())
			} else if err != nil {
				return errres.InternalError(err, c.Logger())
			}

			return c.JSON(http.StatusOK, TicketNewResponse{
				FoundBackfill: true,
				BackfillGsip:  backfill,
			})
		}
	}

	t, ch := tickets.MakeTicket(formData.Class)
	err = engine.AddToPool(t)
	if errs.IsErrorIndexOutOfRange(err) {
		return errres.BadRequest(err, c.Logger())
	} else if err != nil {
		return errres.InternalError(err, c.Logger())
	}

	id := t.Id()
	conn := conns.NewConn(engine.ConnParams(), ch, c.Logger())
	ctx.ConnMap().Set(id, conn)

	return c.JSON(http.StatusOK, TicketNewResponse{
		Id: id,
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

	if conn.Established() {
		return errres.BadRequest(errs.ErrorDuplicatedConnection(), c.Logger())
	}

	ws, err := ctx.WsUpgrader().Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return errres.InternalError(err, c.Logger())
	}

	err = ctx.Engine().PoolToQueue(p.Id)
	if err != nil {
		return errres.BadRequest(err, c.Logger())
	}

	conn.SetWs(ws)
	connMap.Set(p.Id, conn)
	conn.StartWaiting(p.Id, func() {
		connMap.Remove(p.Id)
		c.Logger().Debugf("conn map len: %d", connMap.Count())
	})

	return nil
}
