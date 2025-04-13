package app

import "github.com/gin-gonic/gin"

func New() *app {
	return &app{}
}

type app struct {
	registryPath   string
	rsaPrivKeyPath string
	rsaPubKeyPath  string
	jwtPrivKeyPath string
	jwtPubKeyPath  string
	SQLFiles       []string
	SQLTableStmts  []string
	SQLIndexStmts  []string
	appName        string
	preStartup     func()
	startup        func(c *gin.Engine) error
	onPanic        func(err interface{})
	e              *gin.Engine
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

func (a *app) WithAppName(name string) *app {
	a.appName = name
	return a
}

func (a *app) OnPreStartup(f func()) *app {
	a.preStartup = f
	return a
}

func (a *app) WithSQLFiles(files ...string) *app {
	a.SQLFiles = files
	return a
}

func (a *app) WithSQLTableStatements(stmts ...string) *app {
	a.SQLTableStmts = stmts
	return a
}

func (a *app) WithSQLIndexStatements(stmts ...string) *app {
	a.SQLIndexStmts = stmts
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
