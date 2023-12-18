package queue

import (
	"ssgb-matching/errs"
	"ssgb-matching/matching/tickets"

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

func (q *Q) Enqueue(ticket tickets.Ticket) error {
	return q.inner.Put(ticket)
}

func (q *Q) DequeueN(n int64) ([]tickets.Ticket, error) {
	interfaces, err := q.inner.Get(n)
	if err != nil {
		return nil, err
	}

	len := len(interfaces)
	ts := make([]tickets.Ticket, 0, len)
	for i := 0; i < len; i++ {
		t, ok := interfaces[i].(tickets.Ticket)
		if !ok {
			return nil, errs.ErrorCastFail("ticket")
		}
		ts = append(ts, t)
	}

	return ts, nil
}

func (q *Q) Dequeue() (tickets.Ticket, error) {
	t, err := q.DequeueN(1)
	if err != nil {
		return tickets.Ticket{}, err
	}

	return t[0], nil
}
