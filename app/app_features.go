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
