package engine

import (
	"ssgb-matching/conns"
	"ssgb-matching/errs"
	"ssgb-matching/logger"
	"time"

	"ssgb-matching/matching/matching"
	"ssgb-matching/matching/queue"
	"ssgb-matching/matching/strategy"
	"ssgb-matching/matching/ticket"
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
	params EngineParams
	qMap   map[int64]*queue.Q
	onRoll matching.RollerFunc
	logger logger.Logger
	ticker *time.Ticker
}

func NewEngine(params EngineParams, logger logger.Logger) (*Engine, error) {
	if params.Classes < 1 {
		params.Classes = 1
	}

	f, err := strategy.StrategyFunc(params.Strategy)
	if err != nil {
		return nil, err
	}

	e := &Engine{
		params: params,
		qMap:   make(map[int64]*queue.Q, params.Classes),
		onRoll: f,
		logger: logger,
		ticker: time.NewTicker(time.Millisecond * time.Duration(params.RollingIntervalMil)),
	}

	for i := int64(1); i <= params.Classes; i++ {
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

func (e *Engine) Enqueue(ticket ticket.Ticket) error {
	idx := ticket.Class()
	if err := e.ValidClassIndex(idx); err != nil {
		return err
	}

	if err := e.qMap[idx].Enqueue(ticket); err != nil {
		return err
	}

	e.logger.Debugf("enqueue a ticket to class[%d]", idx)

	return nil
}

func (e *Engine) StartRolling() {
	e.logger.Debugf(
		`engine starts rolling with params: classes->%d, strategy->%d, interval-mil->%d, q-capacity->%d`,
		e.params.Classes,
		e.params.Strategy,
		e.params.RollingIntervalMil,
		e.params.QParams.InitialCapacity,
	)
	go e.roll()
}

func (e *Engine) recover() {
	if r := recover(); r != nil {
		e.logger.Warn("recovering rolling")
		go e.roll()
	}
}

func (e *Engine) roll() {
	defer e.recover()

	for range e.ticker.C {
		_, err := e.onRoll(e.params.MatchingParams, e.qMap)
		if err != nil {
			e.logger.Panic(err)
		}
	}
}
