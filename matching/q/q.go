package q

import (
	"ssgb-matching/errs"
	"ssgb-matching/matching/ticket"

	libqueue "github.com/Workiva/go-datastructures/queue"
)

type QParams struct {
	InitialCapacity int64
}

type Q struct {
	params QParams
	inner  libqueue.Queue
}

func NewQ(params QParams) *Q {
	return &Q{
		params: params,
		inner:  *libqueue.New(params.InitialCapacity),
	}
}

func (q *Q) Len() int64 {
	return q.inner.Len()
}

func (q *Q) Enqueue(ticket ticket.Ticket) error {
	return q.inner.Put(ticket)
}

func (q *Q) DequeueN(n int64) ([]ticket.Ticket, error) {
	interfaces, err := q.inner.Get(n)
	if err != nil {
		return nil, err
	}

	len := len(interfaces)
	tickets := make([]ticket.Ticket, 0, len)
	for i := 0; i < len; i++ {
		t, ok := interfaces[i].(ticket.Ticket)
		if !ok {
			return nil, errs.ErrorCastFail("ticket")
		}
		tickets = append(tickets, t)
	}

	return tickets, nil
}

func (q *Q) Dequeue() (ticket.Ticket, error) {
	t, err := q.DequeueN(1)
	if err != nil {
		return ticket.Ticket{}, err
	}

	return t[0], nil
}
