package app

type featureOpt struct {
	key   string
	value interface{}
}

func LoggingAPI() LoggingAPIFeature {
	return LoggingAPIFeature{
		Enabled: true,
	}
}

type LoggingAPIFeature struct {
	Enabled bool
}

func WithConfig(cfg *AppConfig) Features {
	return Features{
		LoggingAPI: LoggingAPIFeature{
			Enabled: cfg.LoggingAPI.Enabled,
		},
		Docs: DocsFeature{
			Enabled:     cfg.DocsConfig.Enabled,
			DocsPath:    cfg.DocsConfig.DocsDir,
			DocsApiPath: cfg.DocsConfig.DocsApiPath,
		},
		TLS: TLSFeature{
			Enabled:        cfg.TLS.Enabled,
			ServerCertFile: cfg.TLS.CertFile,
			ServerKeyFile:  cfg.TLS.KeyFile,
			CAFile:         cfg.TLS.CaPath,
		},
		JWT: JWTFeature{
			Enabled:                cfg.JWT.Enabled,
			tokenConfigurationPath: cfg.JWT.TokenConfigurationPath,
			PrivateKeyPath:         cfg.JWT.RSAKeyPath,
			PubKeyPath:             cfg.JWT.RSAPubKeyPath,
		},
		SQL: SQLFeature{
			Enabled:               cfg.SQLFiles.Enabled,
			SQLFiles:              cfg.SQLFiles.SQLFiles,
			CreateTableStatements: cfg.SQLFiles.CreateTableStmts,
			CreateIndexStatements: cfg.SQLFiles.CreateIndexStmts,
		},
		Registry: RegistryFeature{
			enabled: cfg.Registry.Enabled,
			registryPath:    &cfg.Registry.Path,
		},
	}
}

type Features struct {
	LoggingAPI LoggingAPIFeature
	RSA        RSAFeature
	JWT        JWTFeature
	SQL        SQLFeature
	HTTP       HTTPClientFeature
	TLS        TLSFeature
	Registry   RegistryFeature
	Docs       DocsFeature
}
