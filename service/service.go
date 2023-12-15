package service

import (
	"ssgb-matching/arg"
	"ssgb-matching/conns"
	"ssgb-matching/matching/engine"
	"ssgb-matching/matching/matching"
	"ssgb-matching/matching/queue"
	"ssgb-matching/server"
	"ssgb-matching/server/context"
	"ssgb-matching/setting"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func Run() {
	e := echo.New()

	args := arg.ParseArgs()
	setting, err := setting.LoadSetting(args.SettingFile)
	if err != nil {
		e.Logger.Fatal(err)
	}

	engine, err := engine.NewEngine(engine.EngineParams{
		Classes:            int64(setting.Classes),
		Strategy:           setting.MatchingStrategy,
		RollingIntervalMil: int64(setting.RollingIntervalMil),
		MatchingParams: matching.MatchingParams{
			MinMatchingCapacity: int64(setting.MinMatchingCapacity),
			MaxMatchingCapacity: int64(setting.MaxMatchingCapacity),
		},
		QParams: queue.QParams{
			InitialCapacity: int64(setting.QInitialCapacity),
		},
		ConnParams: conns.ConnParams{
			ReportIntervalSec: int64(setting.ConnReportIntervalSec),
		},
	}, e.Logger)
	if err != nil {
		e.Logger.Fatal(err)
	}

	c := context.NewComponents(
		context.NewMetadata(setting.ServiceName, setting.ServiceVersion),
		&websocket.Upgrader{},
		engine,
		conns.NewConnMap(),
	)
	s := server.NewServer(
		server.ServerParams{
			ListenAt: setting.ServiceListenAt,
			LogLevel: log.Lvl(setting.LogLevel),
		}, e, c,
	)

	engine.StartRolling()
	e.Logger.Fatal(<-s.Run())
}
