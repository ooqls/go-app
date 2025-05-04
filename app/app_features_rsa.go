package app

const (
	rsa_privateKeyPathOpt string = "opt-private-key-path"
	rsa_publicKeyPathOpt  string = "opt-public-key-path"
)

type rsaOpt = featureOpt

func WithPrivateKeyPath(p string) rsaOpt {
	return rsaOpt{
		key:   rsa_privateKeyPathOpt,
		value: p,
	}
}
func RSA(opts ...rsaOpt) RSAFeature {
	f := RSAFeature{
		Enabled:        true,
		PrivateKeyPath: RsaPrivKeyPathFlag,
		PublicKeyPath:  RsaPubKeyPathFlag,
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

func (f *RSAFeature) apply(opt rsaOpt) {
	switch opt.key {
	case rsa_privateKeyPathOpt:
		f.PrivateKeyPath = opt.value.(string)
	case rsa_publicKeyPathOpt:
		f.PublicKeyPath = opt.value.(string)
	}
}
