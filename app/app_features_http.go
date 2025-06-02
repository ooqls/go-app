package app

import (
	"crypto"
	"crypto/x509"
)

const (
	http_caClientOpt      string = "opt-ca-http-client"
	http_clientCertOpt string = "opt-client-cert-http-client"
	http_clientPrivateKeyOpt  string = "opt-client-key-http-client"
)

type httpOpt struct {
	featureOpt
}

func WithCertificates(certs []x509.Certificate) httpOpt {
	return httpOpt{
		featureOpt: featureOpt{
			key:   http_clientCertOpt,
			value: certs,
		},
	}
}

func WithCaPool(pool x509.CertPool) httpOpt {
	return httpOpt{
		featureOpt: featureOpt{
			key:   http_caClientOpt,
			value: pool,
		},
	}
}

func WithClientCertificates(cert []x509.Certificate) httpOpt {
	return httpOpt{
		featureOpt: featureOpt{
			key:   http_clientCertOpt,
			value: cert,
		},
	}
}

func WithPrivateKey(key crypto.PrivateKey) httpOpt {
	return httpOpt{
		featureOpt: featureOpt{
			key: http_clientPrivateKeyOpt,
			value: key,
		},
	}
}

type HTTPClientFeature struct {
	Enabled            bool
	CA                 *x509.CertPool
	ClientCertificates []x509.Certificate
	PrivateKey         *crypto.PrivateKey
}

func HTTPClient(opts ...httpOpt) HTTPClientFeature {
	f := HTTPClientFeature{
		Enabled: true,
	}

	for _, o := range opts {
		switch o.key {
		case http_caClientOpt:
			f.CA = o.value.(*x509.CertPool)
		case http_clientCertOpt:
			f.ClientCertificates = o.value.([]x509.Certificate)
		case http_clientPrivateKeyOpt:
			f.PrivateKey = o.value.(*crypto.PrivateKey)
		}
	}

	return f
}
