package service

import (
	"ssgb-matching/arg"
	"ssgb-matching/matching/engine"
	"ssgb-matching/matching/q"
	"ssgb-matching/server"
	"ssgb-matching/server/context"
	"ssgb-matching/setting"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func Run() {
	e := echo.New()

	arg := arg.ParseArgs()
	setting, err := setting.LoadSetting(arg.SettingFile)
	if err != nil {
		e.Logger.Fatal(err)
	}

	c := context.NewComponents(
		context.NewMetadata(setting.ServiceName, setting.ServiceVersion),
		&websocket.Upgrader{},
		engine.NewEngine(engine.EngineParams{
			Classes:            setting.Classes,
			RollingIntervalMil: time.Duration(setting.RollingIntervalMil),
			QParams: q.QParams{
				InitialCapacity: int64(setting.QInitialCapacity),
			},
		}, e.Logger),
	)
	s := server.NewServer(
		server.ServerParams{
			ListenAt: setting.ServiceListenAt,
			LogLevel: log.Lvl(setting.LogLevel),
		}, e, c,
	)

	e.Logger.Fatal(<-s.Run())
}
