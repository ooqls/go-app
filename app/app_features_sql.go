package app

import "strings"

type sqlPackage string

const (
	pgxPackage  sqlPackage = "PGX"
	sqlxPackage sqlPackage = "SQLX"
)

const (
	optCreateTableStatements string = "opt-create-table"
	optCreateIndexStatements string = "opt-create-index"
	optSqlFiles string = "opt-sql-files"
)

func WithCreateTableStatements(stmts []string) featureOpt {
	return featureOpt{
		key: optCreateTableStatements,
		value: stmts,
	}
}

func WithCreateIndexStatements(stmts []string) featureOpt {
	return featureOpt{
		key: optCreateIndexStatements,
		value: stmts,
	}
}

func WithSQLFiles(files []string) featureOpt {
	return featureOpt{
		key: optSqlFiles,
		value: files,
	}
}

func SQLXFeature(createTableStatements []string, createIndexStatements []string, sqlFiles []string) SQLFeature {
	return SQLFeature{
		Enabled:               true,
		SQLFiles:              sqlFiles,
		CreateTableStatements: createTableStatements,
		CreateIndexStatements: createIndexStatements,
		SQLPackage:            sqlxPackage,
	}
}

func PGX(createTableStatements, createIndexStatements []string) SQLFeature {
	return SQLFeature{
		Enabled:               true,
		CreateTableStatements: createTableStatements,
		CreateIndexStatements: createIndexStatements,
		SQLFiles:              strings.Split(sqlFiles, ","),
		SQLPackage:            pgxPackage,
	}
}

type SQLFeature struct {
	Enabled               bool
	CreateTableStatements []string
	CreateIndexStatements []string
	SQLFiles              []string
	SQLPackage            sqlPackage
}

func (f *SQLFeature) apply(opt featureOpt) {
	switch opt.key {
	case optCreateIndexStatements:
		f.CreateIndexStatements = opt.value.([]string)
	case optCreateTableStatements:
		f.CreateTableStatements = opt.value.([]string)
	case optSqlFiles:
		f.SQLFiles = opt.value.([]string)
	}
}