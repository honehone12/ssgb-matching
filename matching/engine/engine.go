package engine

import (
	"ssgb-matching/conns"
	"ssgb-matching/errs"
	"ssgb-matching/gsip"
	"ssgb-matching/logger"
	"time"

	"ssgb-matching/matching/matching"
	"ssgb-matching/matching/queue"
	"ssgb-matching/matching/strategy"
	"ssgb-matching/matching/tickets"
)

type EngineParams struct {
	Classes            int64
	Strategy           int
	RollingIntervalMil int64
	MatchingParams     matching.MatchingParams
	QParams            queue.QParams
	ConnParams         conns.ConnParams
}

type Engine struct {
	params     EngineParams
	qMap       map[int64]*queue.Q
	ticketPool *TicketPool
	onRoll     matching.RollerFunc
	provider   gsip.Provider
	logger     logger.Logger
	ticker     *time.Ticker
}

func NewEngine(params EngineParams, provider gsip.Provider, logger logger.Logger,
) (*Engine, error) {
	if params.Classes < 1 {
		params.Classes = 1
	}

	f, err := strategy.StrategyFunc(params.Strategy)
	if err != nil {
		return nil, err
	}

	e := &Engine{
		params:     params,
		qMap:       make(map[int64]*queue.Q, params.Classes),
		ticketPool: NewTicketPool(),
		onRoll:     f,
		provider:   provider,
		logger:     logger,
		ticker:     time.NewTicker(time.Millisecond * time.Duration(params.RollingIntervalMil)),
	}

	for i := int64(1); i <= int64(params.Classes); i++ {
		e.qMap[i] = queue.NewQ(params.QParams)
	}

	return e, nil
}

func (e *Engine) ConnParams() conns.ConnParams {
	return e.params.ConnParams
}

func (e *Engine) ValidClassIndex(idx int64) error {
	if idx < 1 || idx > int64(len(e.qMap)) {
		return errs.ErrorIndexOutOfRange()
	}
	return nil
}

func (e *Engine) AddToPool(ticket tickets.Ticket) error {
	if err := e.ValidClassIndex(ticket.Class()); err != nil {
		return err
	}
	return e.ticketPool.Put(ticket)
}

func (e *Engine) PoolToQueue(id string) error {
	t, err := e.ticketPool.Pull(id)
	if err != nil {
		return err
	}

	idx := t.Class()
	if err := e.qMap[int64(idx)].Enqueue(t); err != nil {
		return err
	}

	e.logger.Debugf("enqueue a ticket to class[%d]", idx)
	return nil
}

func (e *Engine) FindBackfill(class int64) (gsip.GSIP, error) {
	if err := e.ValidClassIndex(class); err != nil {
		return gsip.GSIP{}, err
	}

	return e.provider.BackFillGsip(class)
}

func (e *Engine) StartRolling() {
	e.logger.Infof("engine starts rolling with params: %#v", e.params)
	go e.roll()
}

func (e *Engine) recover() {
	if r := recover(); r != nil {
		e.logger.Warn("recovering rolling")
		go e.roll()
	} else {
		e.ticker.Stop()
		e.logger.Debug("engine stopped rolling")
	}
}

func (e *Engine) roll() {
	defer e.recover()

	for range e.ticker.C {
		matchingResults, err := e.onRoll(e.params.MatchingParams, e.qMap)
		if err != nil {
			e.logger.Panic(err)
		}

		lenResults := len(matchingResults)
		if lenResults == 0 {
			continue
		}

		e.logger.Debugf("%d matching done", lenResults)
		for class, q := range e.qMap {
			e.logger.Debugf("queue[%d] length: %d", class, q.Len())
		}

		for i := 0; i < lenResults; i++ {
			matched := matchingResults[i]
			lenTickets := len(matched.Matched)
			result := gsip.GSIPResult{}

			ip, err := e.provider.NextGsip(matched.Class)
			if err != nil {
				e.logger.Error(err)
				result.Status = gsip.StatusError
			} else {
				result.Status = gsip.StatusOk
				result.Gsip = ip
			}

			for j := 0; j < lenTickets; j++ {
				matched.Matched[j].Chan() <- result
			}
		}
	}
}
