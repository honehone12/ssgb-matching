package context

import (
	"ssgb-matching/matching/engine"

	"github.com/gorilla/websocket"
)

type Components struct {
	metadata   *Metadata
	wsUpgrader *websocket.Upgrader
	engine     *engine.Engine
}

func NewComponents(
	metadata *Metadata,
	wsUpgrader *websocket.Upgrader,
	engine *engine.Engine,
) *Components {
	return &Components{
		metadata:   metadata,
		wsUpgrader: wsUpgrader,
		engine:     engine,
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
