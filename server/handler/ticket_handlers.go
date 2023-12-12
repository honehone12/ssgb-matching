package handler

import (
	"net/http"
	"ssgb-matching/errs"
	"ssgb-matching/matching/ticket"
	"ssgb-matching/server/context"
	"ssgb-matching/server/errres"
	"ssgb-matching/server/form"

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

	t := ticket.MakeTicket(formData.Class)
	err = ctx.Components.Engine().Enqueue(t)
	if errs.IsErrorIndexOutOfRange(err) {
		return errres.BadRequest(err, c.Logger())
	} else if err != nil {
		return errres.InternalError(err, c.Logger())
	}

	return c.JSON(http.StatusOK, TicketNewResponse{
		Id: t.Id(),
	})
}
