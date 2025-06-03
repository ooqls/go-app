package app

import (
	"context"
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ooqls/go-log"
	"go.uber.org/zap"
)

var registryPathFlag string

var sqlFilesFlag string

var RsaPrivKeyPathFlag string
var RsaPubKeyPathFlag string

var JwtPrivKeyPathFlag string
var JwtPubKeyPathFlag string

var caBundlePathFlag string

func init() {
	flag.StringVar(&registryPathFlag, "registry", "", "Path to the registry path")
	flag.StringVar(&sqlFilesFlag, "sql-files", "", "Comma separated list of files")
	flag.StringVar(&RsaPrivKeyPathFlag, "rsa-private-key", "", "Path to an RSA private key")
	flag.StringVar(&RsaPubKeyPathFlag, "rsa-public-key", "", "Path to the RSA public key")
	flag.StringVar(&JwtPrivKeyPathFlag, "jwt-private-key", "", "Path to a JWT private key")
	flag.StringVar(&JwtPubKeyPathFlag, "jwt-public-key", "", "Path to a jwt public key")
	flag.StringVar(&caBundlePathFlag, "ca-bundle", "", "Path to a ca bundle")
}

func New(appName string, features Features) *app {
	return &app{
		appName:  appName,
		e:        gin.New(),
		l:        log.NewLogger(appName),
		features: features,
	}
}

type app struct {
	appName    string
	preStartup func()
	startup    func(ctx *StartupContext) error
	onPanic    func(err interface{})
	e          *gin.Engine
	l          *zap.Logger
	state      AppState
	features   Features
	testEnvironment *TestEnvironment
}

func (a *app) WithTestEnvironment(env TestEnvironment) {
	a.testEnvironment = &env
}


func (a *app) OnPreStartup(f func()) *app {
	a.preStartup = f
	return a
}

func (a *app) OnStartup(f func(ctx *StartupContext) error) *app {
	a.startup = f
	return a
}

func (a *app) Run() error {
	flag.Parse()
	if a.testEnvironment != nil {
		cleanup, err := a.testEnvironment.Start(context.Background())
		if err != nil {
			return fmt.Errorf("failed to start test environment: %v", err)
		}
		defer cleanup()
	}

	if err := a._startup(); err != nil {
		return err
	}

	return nil
}
