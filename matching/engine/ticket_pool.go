package engine

import (
	"ssgb-matching/errs"
	"ssgb-matching/logger"
	"ssgb-matching/matching/tickets"
	"sync"
	"time"
)

type TicketPoolParams struct {
	CleanUpIntervalSec int64
	WsUpgradeLimitSec  int64
}

type poolItem struct {
	ticket     tickets.Ticket
	timePooled time.Time
}

type TicketPool struct {
	params  TicketPoolParams
	count   int64
	inner   sync.Map
	ticker  time.Ticker
	closeCh chan bool
	logger  logger.Logger
}

func NewTicketPool(params TicketPoolParams, logger logger.Logger) *TicketPool {
	return &TicketPool{
		params:  params,
		count:   0,
		inner:   sync.Map{},
		ticker:  *time.NewTicker(time.Second),
		closeCh: make(chan bool),
		logger:  logger,
	}
}

func (p *TicketPool) recover() {
	if r := recover(); r != nil {
		p.logger.Warn("recovering clean up")
		go p.CleanUp()
	}
}

func (p *TicketPool) CleanUp() {
	defer p.recover()

LOOP:
	for {
		select {
		case <-p.closeCh:
			p.ticker.Stop()
			break LOOP
		case <-p.ticker.C:
			tmp := make([]string, 0)
			limit := time.Now().Add(-time.Second * time.Duration(p.params.WsUpgradeLimitSec))
			p.logger.Debugf("cleaning up pool count: %d", p.count)
			p.inner.Range(func(k, v interface{}) bool {
				item, ok := v.(poolItem)
				if !ok {
					p.logger.Panic(errs.ErrorCastFail("poolitem"))
				}

				if item.timePooled.Before(limit) {
					id, ok := k.(string)
					if !ok {
						p.logger.Panic(errs.ErrorCastFail("string"))
					}

					tmp = append(tmp, id)
				}
				return true
			})

			for i, count := 0, len(tmp); i < count; i++ {
				id := tmp[i]
				p.inner.Delete(id)
				p.count--
				p.logger.Warnf("removed from pool [%s]", id)
			}
		}
	}
}

func (p *TicketPool) CloseCh() chan<- bool {
	return p.closeCh
}

func (p *TicketPool) Count() int64 {
	return p.count
}

func (p *TicketPool) Put(t tickets.Ticket) error {
	_, exists := p.inner.LoadOrStore(t.Id(), poolItem{
		ticket:     t,
		timePooled: time.Now(),
	})
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
	item, ok := i.(poolItem)
	if !ok {
		return tickets.Ticket{}, errs.ErrorCastFail("ticket")
	}

	return item.ticket, nil
}
