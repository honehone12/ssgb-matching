package tickets

import (
	"ssgb-matching/gsip"

	libuuid "github.com/google/uuid"
)

type Ticket struct {
	id     string
	class  int64
	gsipCh chan<- gsip.GSIP
}

func MakeTicket(class int64) (Ticket, <-chan gsip.GSIP) {
	ch := make(chan gsip.GSIP)
	return Ticket{
		id:     libuuid.NewString(),
		class:  class,
		gsipCh: ch,
	}, ch
}

func (t *Ticket) Id() string {
	return t.id
}

func (t *Ticket) Class() int64 {
	return t.class
}

func (t *Ticket) Chan() chan<- gsip.GSIP {
	return t.gsipCh
}
