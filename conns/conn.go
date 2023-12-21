package conns

import (
	"ssgb-matching/gsip"
	"ssgb-matching/logger"
	"ssgb-matching/messages"
	"time"

	"github.com/gorilla/websocket"
)

type ConnParams struct {
	ReportIntervalSec int64
	WsTimeoutSec      int64
}

type Conn struct {
	params ConnParams
	inner  *websocket.Conn
	ticker *time.Ticker
	waitCh <-chan gsip.GSIPResult
	logger logger.Logger
}

func NewConn(params ConnParams, waitCh <-chan gsip.GSIPResult, logger logger.Logger) *Conn {
	return &Conn{
		params: params,
		ticker: time.NewTicker(time.Second * time.Duration(params.ReportIntervalSec)),
		waitCh: waitCh,
		logger: logger,
	}
}

func (c *Conn) Established() bool {
	return c.inner != nil
}

func (c *Conn) SetWs(conn *websocket.Conn) {
	c.inner = conn
}

func (c *Conn) StartWaiting(id string, onDone func()) {
	go c.sendMessagedWhileWait(id, onDone)
}

func (c *Conn) recover(id string, onDone func()) {
	if r := recover(); r != nil {
		c.logger.Warn("recovering sending")
		go c.sendMessagedWhileWait(id, onDone)
	} else {
		c.ticker.Stop()
		if c.Established() {
			c.inner.Close()
		}
		onDone()
		c.logger.Debugf("conn[%s] closed", id)
	}
}

func (c *Conn) sendMessagedWhileWait(id string, onDone func()) {
	defer c.recover(id, onDone)

WAIT:
	for {
		msg := messages.StatusMessage{}
		select {
		case result := <-c.waitCh:
			if !c.Established() {
				c.logger.Warnf("conn[%s] matched, but is already broken", id)
				break WAIT
			}

			switch result.Status {
			case gsip.StatusOk:
				msg.Status = messages.StatusMatched
				msg.Gsip = result.Gsip
			default:
				msg.Status = messages.StatusError
			}
		case <-c.ticker.C:
			if !c.Established() {
				continue
			}

			msg.Status = messages.StatusWaitng
		}

		timeout := time.Now().Add(time.Second * time.Duration(c.params.WsTimeoutSec))
		if err := c.inner.SetWriteDeadline(timeout); err != nil {
			c.logger.Panic(err)
		}

		if err := c.inner.WriteJSON(msg); err != nil {
			c.logger.Warnf("conn[%s] time out: %s", id, err)
			c.inner.Close()
			c.SetWs(nil)
		}

		if msg.Status == messages.StatusMatched || msg.Status == messages.StatusError {
			break
		}
	}
}
