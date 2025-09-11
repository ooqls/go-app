package app

type featureOpt struct {
	key   string
	value interface{}
}

func WithConfig(cfg *AppConfig) Features {
	return Features{
		LoggingAPI: LoggingApiFeature{
			Enabled: cfg.LoggingAPI.Enabled,
			Port:    cfg.LoggingAPI.Port,
		},
		Gin: GinFeature{
			Enabled: cfg.Gin.Enabled,
			Port:    cfg.Gin.Port,
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
			Enabled:                 cfg.JWT.Enabled,
			tokenConfigurationPaths: cfg.JWT.TokenConfigurationPaths,
			PrivateKeyPath:          cfg.JWT.RSAKeyPath,
			PubKeyPath:              cfg.JWT.RSAPubKeyPath,
			tokenConfiguration:      cfg.JWT.TokenConfigurations,
		},
		Health: HealthFeature{
			Enabled:  cfg.Health.Enabled,
			Path:     cfg.Health.Path,
			Interval: cfg.Health.Interval,
		},
		SQL: SQLFeature{
			Enabled:               cfg.SQLFiles.Enabled,
			SQLPackage:            cfg.SQLFiles.SQLPackage,
			SQLFiles:              cfg.SQLFiles.SQLFiles,
			SQLDirs:               cfg.SQLFiles.SQLFilesDirs,
			CreateTableStatements: cfg.SQLFiles.CreateTableStmts,
			CreateIndexStatements: cfg.SQLFiles.CreateIndexStmts,
		},
		Registry: RegistryFeature{
			enabled:      cfg.Registry.Enabled,
			registryPath: &cfg.Registry.Path,
		},
	}
}

type Features struct {
	LoggingAPI LoggingApiFeature
	RSA        RSAFeature
	JWT        JWTFeature
	SQL        SQLFeature
	HTTP       HTTPFeature
	TLS        TLSFeature
	Registry   RegistryFeature
	Docs       DocsFeature
	Health     HealthFeature
	Gin        GinFeature
}
