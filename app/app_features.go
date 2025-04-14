package app

type featureOpt struct {
	key   string
	value interface{}
}

const (
	optPrivateKeyPath string = "opt-private-key-path"
	optPublicKeyPath string = "opt-public-key-path"

	
)

func WithPrivateKeyPath(p string) featureOpt {
	return featureOpt{
		key: optPrivateKeyPath,
		value: p,
	}
}
func RSA(opts ...featureOpt) RSAFeature {
	f := RSAFeature{
		Enabled:        true,
		PrivateKeyPath: RsaPrivKeyPath,
		PublicKeyPath:  RsaPubKeyPath,
	}
	for _, opt := range opts {
		f.apply(opt)
	}

	return f
}

type RSAFeature struct {
	Enabled        bool
	PrivateKeyPath string
	PublicKeyPath  string
}

func (f *RSAFeature) apply(opt featureOpt) {
	switch opt.key {
	case optPrivateKeyPath:
		f.PrivateKeyPath = opt.value.(string)
	case optPublicKeyPath:
		f.PublicKeyPath = opt.value.(string)
	}
}

func JWT(opts ...featureOpt) JWTFeature {
	f := JWTFeature{
		Enabled:        true,
		PrivateKeyPath: JwtPrivKeyPath,
		PubKeyPath:     JwtPubKeyPath,
	}

	for _, opt := range opts {
		f.apply(opt)
	}

	return f
}

type JWTFeature struct {
	Enabled        bool
	PrivateKeyPath string
	PubKeyPath     string
}

func (f *JWTFeature) apply(opt featureOpt) {
	switch opt.key {
	case optPrivateKeyPath:
		f.PrivateKeyPath = opt.value.(string)
	case optPublicKeyPath:
		f.PubKeyPath = opt.value.(string)
	}
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
}
