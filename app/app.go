package app

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/ooqls/go-log"
	"go.uber.org/zap"
)

func New(appName string, features Features) *app {
	flag.Parse()
	
	return &app{
		appName:        appName,
		e:              gin.New(),
		l:              log.NewLogger(appName),
	}
}

type app struct {
	appName        string
	preStartup     func()
	startup        func(c *gin.Engine) error
	onPanic        func(err interface{})
	e              *gin.Engine
	l              *zap.Logger

	state    AppState
	features Features
}

func (a *app) OnPreStartup(f func()) *app {
	a.preStartup = f
	return a
}

func (a *app) OnStartup(f func(e *gin.Engine) error) *app {
	a.startup = f
	return a
}

func (a *app) Run() error {
	if err := a._startup(); err != nil {
		return err
	}

	return nil
}
