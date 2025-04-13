package app

import (
	"github.com/jmoiron/sqlx"
	"github.com/ooqls/go-crypto/keys"
	gosqlx "github.com/ooqls/go-db/sqlx"
	"github.com/ooqls/go-log"
	v1 "github.com/ooqls/go-log/api/v1"
	"github.com/ooqls/go-registry"
	"go.uber.org/zap"
)

func (a *app) _startup_rsa() error {
	l := log.NewLogger(a.appName)
	if a.rsaPrivKeyPath != "" && a.rsaPubKeyPath != "" {
		l.Debug("[Startup RSA] initializing RSA keys with paths",
			zap.String("rsa_private_key_path", a.rsaPrivKeyPath),
			zap.String("rsa_pub_key_path", a.rsaPubKeyPath),
		)
		if !fileExists(a.rsaPrivKeyPath) {
			l.Info("[Startup RSA] RSA private key not found", zap.String("path", a.rsaPrivKeyPath))
			return nil
		}

		if !fileExists(a.rsaPubKeyPath) {
			l.Info("[Startup RSA] RSA public key not found", zap.String("path", a.rsaPubKeyPath))
			return nil
		}
		privKey := mustReadFile(a.rsaPrivKeyPath)
		pubKey := mustReadFile(a.rsaPubKeyPath)

		if err := keys.InitRSA(privKey, pubKey); err != nil {
			return err
		}
		l.Debug("[Startup RSA] RSA keys initialized successfully")
	} else {
		l.Debug("[Startup RSA] no RSA key paths provided")
	}

	return nil
}

func (a *app) _startup_jwt() error {
	l := a.l
	if a.jwtPrivKeyPath != "" && a.jwtPubKeyPath != "" {
		l.Debug("[Startup JWT] initializing JWT keys with paths",
			zap.String("jwt_private_key_path", a.jwtPrivKeyPath),
			zap.String("jwt_pub_key_path", a.jwtPubKeyPath),
		)
		jwtPrivKey := mustReadFile(a.jwtPrivKeyPath)
		jwtPubKey := mustReadFile(a.jwtPubKeyPath)
		if err := keys.InitJwt(jwtPrivKey, jwtPubKey); err != nil {
			return err
		}

		l.Debug("[Startup JWT] JWT keys initialized successfully")
	} else {
		l.Debug("[Startup JWT] no JWT key paths provided")
	}

	return nil
}

func (a *app) _startup_registry() error {
	l := log.NewLogger(a.appName)

	if a.registryPath != "" {
		if !fileExists(a.registryPath) {
			l.Info("[Startup Registry] registry file not found", zap.String("registry_path", a.registryPath))
			return nil
		}

		l.Debug("[Startup Registry] initializing registry with path", zap.String("path", a.registryPath))
		if err := registry.Init(a.registryPath); err != nil {
			l.Error("[Startup Registry] failed to initialize registry", zap.Error(err))
			return err
		}
		l.Debug("[Startup Registry] registry initialized successfully")

	} else {
		l.Debug("[Startup Registry] no registry path provided")
	}

	return nil

}

func (a *app) _startup_sql() error {
	l := a.l
	if len(a.SQLFiles) > 0 {
		l.Debug("[Startup SQL] initializing SQL files", zap.Strings("sql_files", a.SQLFiles))
		err := gosqlx.InitDefault()
		if err != nil {
			l.Error("[Startup SQL] failed to initialize SQL files", zap.Error(err))
			return err
		}

		c := gosqlx.GetSQLX()
		for _, file := range a.SQLFiles {
			l.Debug("[Startup SQL] Loading sql file: " + file)
			if _, err := sqlx.LoadFile(c, file); err != nil {
				l.Error("[Startup SQL] failed to load file: "+file, zap.Error(err))
			}
		}
		l.Debug("[Startup SQL] SQL files initialized successfully")
	}

	tableStmts := []string{}
	indexStmts := []string{}

	if len(a.SQLTableStmts) > 0 {
		tableStmts = a.SQLTableStmts
	}

	if len(a.SQLIndexStmts) > 0 {
		indexStmts = a.SQLIndexStmts
	}

	if len(indexStmts) > 0 || len(tableStmts) > 0 {
		l.Debug("[Startup SQL] seeding with SQL statements")
		gosqlx.SeedSQLX(tableStmts, indexStmts)
		l.Debug("[Startup SQL] finished seeding with SQL statements")
	} else {
		l.Debug("[Startup SQL] no SQL statmenets")
	}

	return nil
}

func (a *app) _startup() error {
	l := log.NewLogger(a.appName).With(zap.String("stage", "startup"))
	if a.onPanic != nil {
		defer func() {
			if err := recover(); err != nil {
				l.Warn("[Startup] recovered from panic", zap.Any("error", err))
				a.onPanic(err)
			}
		}()
	}

	if a.preStartup != nil {
		l.Debug("[Startup] running pre-startup function")
		a.preStartup()
		l.Debug("[Startup] pre-startup function completed")
	}

	startup_funcs := []func()error{
		a._startup_registry,
		a._startup_rsa,
		a._startup_jwt,
		a._startup_sql,
	}

	for _, f := range startup_funcs {
		err := f()
		if err != nil {
			return err
		}
	}


	l.Debug("[Startup] dding logging routes")
	v1.AddRoutes(a.e)
	l.Debug("[Startup] finished adding logging routes")
	if a.startup != nil {
		l.Debug("[Startup] running startup function")
		err := a.startup(a.e)
		if err != nil {
			l.Error("encountered an error on startup", zap.Error(err))
			return err
		}
		l.Debug("[Startup] startup function completed")

	}

	return nil
}
