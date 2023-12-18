package engine

import (
	"ssgb-matching/errs"
	"ssgb-matching/matching/tickets"
	"sync"
)

type TicketPool struct {
	count int64
	inner sync.Map
}

func NewTicketPool() *TicketPool {
	return &TicketPool{
		count: 0,
		inner: sync.Map{},
	}
}

func (p *TicketPool) Len() int64 {
	return p.count
}

func (p *TicketPool) Put(t tickets.Ticket) error {
	_, exists := p.inner.LoadOrStore(t.Id(), t)
	if exists {
		return errs.ErrorItemAlreadyExists()
	}

	p.count++
	return nil
}

func (p *TicketPool) Pull(id string) (tickets.Ticket, error) {
	i, exists := p.inner.LoadAndDelete(id)
	if !exists {
		return tickets.Ticket{}, errs.ErrorNoSuchItem()
	}

	p.count--
	t, ok := i.(tickets.Ticket)
	if !ok {
		return tickets.Ticket{}, errs.ErrorCastFail("ticket")
	}

	return t, nil
}
