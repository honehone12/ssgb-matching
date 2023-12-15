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
	params  ConnParams
	inner   *websocket.Conn
	ticker  *time.Ticker
	waitCh  <-chan gsip.GSIP
	closeCh chan gsip.GSIP
	logger  logger.Logger
}

func MakeConn(params ConnParams, waitCh <-chan gsip.GSIP, logger logger.Logger) Conn {
	return Conn{
		params:  params,
		ticker:  time.NewTicker(time.Second * time.Duration(params.ReportIntervalSec)),
		waitCh:  waitCh,
		closeCh: make(chan gsip.GSIP),
		logger:  logger,
	}
}

func (c *Conn) Established() bool {
	return c.inner != nil
}

func (c *Conn) SetWs(conn *websocket.Conn) {
	c.inner = conn
}

func (c *Conn) wait() {
	c.closeCh <- <-c.waitCh
}

func (c *Conn) recover(id string) {
	if r := recover(); r != nil {
		c.logger.Warn("recovering sending")
		go c.sendWhileWait(id)
	}
}

func (c *Conn) sendWhileWait(id string) {
	if !c.Established() {
		return
	}

	defer c.recover(id)

	for {
		var msg messages.StatusMessage
		select {
		case gsip := <-c.closeCh:
			msg = messages.MakeMatchedMessage(gsip)
		case <-c.ticker.C:
			msg = messages.MakeWaitingMessage()
		}

		timeout := time.Now().Add(time.Duration(c.params.WsTimeoutSec))
		if err := c.inner.SetWriteDeadline(timeout); err != nil {
			c.logger.Panic(err)
		}

		if err := c.inner.WriteJSON(msg); err != nil {
			c.logger.Warnf("conn[%s] time out: %s", id, err)
		}

		if msg.Status == messages.StatusMatched {
			break
		}
	}
}
