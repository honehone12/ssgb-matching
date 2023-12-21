package tickets

import (
	"ssgb-matching/gsip"

	libuuid "github.com/google/uuid"
)

type Ticket struct {
	id     string
	class  int64
	gsipCh chan<- gsip.GSIPResult
}

func MakeTicket(class int64) (Ticket, <-chan gsip.GSIPResult) {
	ch := make(chan gsip.GSIPResult)
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

func (t *Ticket) Chan() chan<- gsip.GSIPResult {
	return t.gsipCh
}
