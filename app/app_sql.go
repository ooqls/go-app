package app

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/ooqls/go-db/pgx"
	gosqlx "github.com/ooqls/go-db/sqlx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (a *app) _seed_pgx_files(ctx *StartupContext) bool {
	l := ctx.L()
	sqlSeeded := len(a.features.SQL.SQLFiles) > 0

	for _, file := range a.features.SQL.SQLFiles {
		l.Debug("loading SQL file: " + file)
		if err := pgx.SeedPGXFile(ctx, file); err != nil {
			l.Error("failed to load file: "+file, zap.Error(err))
			sqlSeeded = false
		}
	}

	return sqlSeeded
}

func (a *app) _seed_sqlx_files(ctx *StartupContext) bool {
	l := ctx.L()
	c := gosqlx.GetSQLX()
	sqlSeeded := len(a.features.SQL.SQLFiles) > 0
	for _, file := range a.features.SQL.SQLFiles {
		l.Debug("Loading sql file: " + file)
		if _, err := sqlx.LoadFile(c, file); err != nil {
			l.Error("failed to load file: "+file, zap.Error(err))
			sqlSeeded = false
		}
	}

	return sqlSeeded
}

func (a *app) _startup_sql(ctx *StartupContext) error {
	l := ctx.L().WithOptions(zap.Hooks(func(e zapcore.Entry) error {
		e.Message = fmt.Sprintf("[Startup %s] %s", a.features.SQL.SQLPackage, e.Message)
		return nil
	}))
	ctx = NewStartupContext(ctx, l)

	if len(a.features.SQL.SQLFiles) > 0 {
		if !a.state.RegistryInitialized {
			l.Info("Please intialize registry before starting SQL")
			return nil
		}

		l.Debug("initializing SQL files", zap.Strings("sql_files", a.features.SQL.SQLFiles))

		if a.features.SQL.SQLPackage == sqlxPackage {
			err := gosqlx.InitDefault()
			if err != nil {
				l.Error("failed to initialize SQLX", zap.Error(err))
				return err
			}

			a.state.SQLSeeded = a._seed_sqlx_files(ctx)
		} else if a.features.SQL.SQLPackage == pgxPackage {
			err := pgx.InitDefault()
			if err != nil {
				l.Error("failed to initialize PGX", zap.Error(err))
				return err
			}

			a.state.SQLSeeded = a._seed_pgx_files(ctx)
		}
		l.Debug("SQL files initialized successfully")
	}

	tableStmts := []string{}
	indexStmts := []string{}

	if len(a.features.SQL.CreateTableStatements) > 0 {
		tableStmts = a.features.SQL.CreateTableStatements
	}

	if len(a.features.SQL.CreateIndexStatements) > 0 {
		indexStmts = a.features.SQL.CreateIndexStatements
	}

	if len(indexStmts) > 0 || len(tableStmts) > 0 {
		l.Debug("seeding with SQL statements")
		if a.features.SQL.SQLPackage == sqlxPackage {
			gosqlx.SeedSQLX(tableStmts, indexStmts)
		} else if a.features.SQL.SQLPackage == pgxPackage {
			pgx.SeedPGX(ctx, tableStmts, indexStmts)
		}

		l.Debug("finished seeding with SQL statements")
	} else {
		l.Debug("no SQL statmenets")
	}

	return nil
}
