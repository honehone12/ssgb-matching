package ticket

import libuuid "github.com/google/uuid"

type Ticket struct {
	id string

	class int64
}

func MakeTicket(class int64) Ticket {
	return Ticket{
		id:    libuuid.NewString(),
		class: class,
	}
}

func (t *Ticket) Id() string {
	return t.id
}

func (t *Ticket) Class() int64 {
	return t.class
}
