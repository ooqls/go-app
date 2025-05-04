package app

import "github.com/ooqls/go-crypto/jwt"

type jwtOpt = featureOpt

const (
	jwt_tokenConfigurationPathOpt string = "jwt_tokenConfigurationPath"
	jwt_tokenConfigurationOpt     string = "jwt_tokenConfiguration"
)

func WithTokenConfigurationPath(p string) featureOpt {
	return featureOpt{
		key: jwt_tokenConfigurationPathOpt,
		value: p,
	}
}

func WithTokenConfiguration(cfg jwt.TokenConfiguration) featureOpt {
	return featureOpt{
		key: jwt_tokenConfigurationOpt,
		value: &cfg,
	}
}

func JWT(opts ...jwtOpt) JWTFeature {
	f := JWTFeature{
		Enabled:        true,
		PrivateKeyPath: JwtPrivKeyPathFlag,
		PubKeyPath:     JwtPubKeyPathFlag,
	}

	for _, opt := range opts {
		f.apply(opt)
	}

	return f
}

type JWTFeature struct {
	Enabled                bool
	PrivateKeyPath         string
	PubKeyPath             string
	tokenConfigurationPath string
	tokenConfiguration     *jwt.TokenConfiguration
}

func (f *JWTFeature) apply(opt jwtOpt) {
	switch opt.key {
	case rsa_privateKeyPathOpt:
		f.PrivateKeyPath = opt.value.(string)
	case rsa_publicKeyPathOpt:
		f.PubKeyPath = opt.value.(string)
	case jwt_tokenConfigurationPathOpt:
		f.tokenConfigurationPath = opt.value.(string)
	case jwt_tokenConfigurationOpt:
		f.tokenConfiguration = opt.value.(*jwt.TokenConfiguration)
	}
}
