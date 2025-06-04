package app

import "github.com/ooqls/go-crypto/jwt"

type jwtOpt struct {
	featureOpt
}

const (
	jwt_tokenConfigurationPathOpt string = "jwt_tokenConfigurationPath"
	jwt_tokenConfigurationOpt     string = "jwt_tokenConfiguration"
)

func WithTokenConfigurationPath(p string) jwtOpt {
	return jwtOpt{
		featureOpt: featureOpt{
			key: jwt_tokenConfigurationPathOpt,
			value: p,
		},
	}
}

func WithTokenConfiguration(cfg jwt.TokenConfiguration) jwtOpt {
	return jwtOpt{
		featureOpt: featureOpt{
			key: jwt_tokenConfigurationOpt,
			value: &cfg,
		},
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
	tokenConfigurationPaths []string
	tokenConfiguration     []jwt.TokenConfiguration
}

func (f *JWTFeature) apply(opt jwtOpt) {
	switch opt.key {
	case rsa_privateKeyPathOpt:
		f.PrivateKeyPath = opt.value.(string)
	case rsa_publicKeyPathOpt:
		f.PubKeyPath = opt.value.(string)
	case jwt_tokenConfigurationPathOpt:
		f.tokenConfigurationPaths = opt.value.([]string)
	case jwt_tokenConfigurationOpt:
		f.tokenConfiguration = opt.value.([]jwt.TokenConfiguration)
	}
}
