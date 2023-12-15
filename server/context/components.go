package context

import (
	"ssgb-matching/conns"
	"ssgb-matching/matching/engine"

	"github.com/gorilla/websocket"
)

type Components struct {
	metadata   *Metadata
	wsUpgrader *websocket.Upgrader
	engine     *engine.Engine
	connMap    *conns.ConnMap
}

func NewComponents(
	metadata *Metadata,
	wsUpgrader *websocket.Upgrader,
	engine *engine.Engine,
	connMap *conns.ConnMap,
) *Components {
	return &Components{
		metadata:   metadata,
		wsUpgrader: wsUpgrader,
		engine:     engine,
		connMap:    connMap,
	}
}

func (c *Components) Metadata() *Metadata {
	return c.metadata
}

func (c *Components) WsUpgrader() *websocket.Upgrader {
	return c.wsUpgrader
}

func (c *Components) Engine() *engine.Engine {
	return c.engine
}

func (c *Components) ConnMap() *conns.ConnMap {
	return c.connMap
}
