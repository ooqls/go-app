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
		registryPath:   registryPath,
		e:              gin.New(),
		l:              log.NewLogger(appName),
	}
}

type app struct {
	registryPath   string
	rsaPrivKeyPath string
	rsaPubKeyPath  string
	jwtPrivKeyPath string
	jwtPubKeyPath  string
	appName        string
	preStartup     func()
	startup        func(c *gin.Engine) error
	onPanic        func(err interface{})
	e              *gin.Engine
	l              *zap.Logger

	state    AppState
	features Features
}

func (a *app) WithRegistryPath(path string) *app {
	a.registryPath = path
	return a
}

func (a *app) WithRsaPath(privKeyPath, pubKeyPath string) *app {
	a.rsaPrivKeyPath = privKeyPath
	a.rsaPubKeyPath = pubKeyPath
	return a
}

func (a *app) WithJwtPath(privKeyPath, pubKeyPath string) *app {
	a.jwtPrivKeyPath = privKeyPath
	a.jwtPubKeyPath = pubKeyPath
	return a
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
