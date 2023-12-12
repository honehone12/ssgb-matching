package engine

import (
	"ssgb-matching/errs"
	"ssgb-matching/logger"
	"time"

	"ssgb-matching/matching/q"
	"ssgb-matching/matching/ticket"
)

type EngineParams struct {
	Classes            int
	RollingIntervalMil time.Duration
	QParams            q.QParams
}

type Engine struct {
	qList  []*q.Q
	logger logger.Logger
}

func NewEngine(params EngineParams, logger logger.Logger) *Engine {
	e := &Engine{
		qList:  make([]*q.Q, 0, params.Classes),
		logger: logger,
	}

	for i := 0; i < params.Classes; i++ {
		e.qList = append(e.qList, q.NewQ(params.QParams))
	}

	return e
}

func classIndex(class int64) int {
	return int(class - 1)
}

func (e *Engine) validClassIndex(idx int) error {
	if idx < 0 || idx >= len(e.qList) {
		return errs.ErrorIndexOutOfRange()
	}
	return nil
}

func (e *Engine) Enqueue(ticket ticket.Ticket) error {
	idx := classIndex(ticket.Class())
	if err := e.validClassIndex(idx); err != nil {
		return err
	}

	return e.qList[idx].Enqueue(ticket)
}
