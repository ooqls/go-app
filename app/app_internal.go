package app

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ooqls/go-crypto/jwt"
	"github.com/ooqls/go-crypto/keys"
	v1 "github.com/ooqls/go-log/api/v1"
	"github.com/ooqls/go-registry"
	"go.uber.org/zap"
)

func (a *app) _startup_docs(ctx *StartupContext) error {
	l := ctx.L()
	l.Info("[Startup docs] Serving htnl docs",
		zap.String("path", a.features.Docs.DocsPath), zap.String("api_path", a.features.Docs.DocsApiPath))
	a.e.GET(a.features.Docs.DocsApiPath, func(ctx *gin.Context) {
		ctx.File(a.features.Docs.DocsPath)
	})

	return nil
}

func (a *app) _startup_http(ctx *StartupContext) error {
	ca := a.features.HTTP.CA
	clientCert := a.features.HTTP.ClientCertificates
	privateKey := a.features.HTTP.PrivateKey

	certChain := make([][]byte, 512)
	for _, crt := range clientCert {
		certChain = append(certChain, crt.Raw)
	}

	tlsCert := tls.Certificate{
		Certificate: certChain,
		PrivateKey:  privateKey,
	}

	trans := http.DefaultTransport.(*http.Transport)
	trans.TLSClientConfig.Certificates = append(trans.TLSClientConfig.Certificates, tlsCert)
	trans.TLSClientConfig.ClientCAs = ca
	trans.TLSClientConfig.RootCAs = ca

	return nil
}

func (a *app) _startup_tls(ctx *StartupContext) error {
	f := a.features.TLS
	l := ctx.L()
	cfg := &tls.Config{}

	pool, err := x509.SystemCertPool()
	if err != nil {
		return fmt.Errorf("failed to get system cert pool: %v", err)
	}

	if len(f.CABytes) > 0 {
		l.Info("[Startup TLS] using CA bytes in tls config")
		ok := pool.AppendCertsFromPEM(f.CABytes)
		if !ok {
			return fmt.Errorf("failed to append all ca bytes")
		}
	} else if f.CAFile != "" {
		l.Info("[Startup TLS] Using CA file for tls config", zap.String("ca", f.CAFile))
		b, err := os.ReadFile(f.CAFile)
		if err != nil {
			return fmt.Errorf("failed to read ca file %s: %v", f.CAFile, err)
		}

		ok := pool.AppendCertsFromPEM(b)
		if !ok {
			return fmt.Errorf("failed to add pem bytes from ca file %s", f.CAFile)
		}
	} else {
		l.Info("[Startup TLS] No CA given")
	}
	cfg.ClientCAs = pool
	cfg.RootCAs = pool

	var tlsCert tls.Certificate
	var keyBytes []byte
	var certBytes []byte

	if len(f.ServerCertBytes) > 0 {
		l.Info("[Startup TLS] server cert bytes given")
		certBytes = f.ServerCertBytes
	} else if f.ServerCertFile != "" {
		l.Info("[Startup TLS] server cert file given", zap.String("cert_file", f.ServerCertFile))
		certBytes, err = os.ReadFile(f.ServerCertFile)
		if err != nil {
			return fmt.Errorf("failed to read cert file %s: %v", f.ServerCertFile, err)
		}
	} else {
		return fmt.Errorf("no server cert given for TLS")
	}

	if len(f.ServerKeyBytes) > 0 {
		l.Info("[Startup TLS] server key bytes given")
		keyBytes = f.ServerKeyBytes
	} else if f.ServerKeyFile != "" {
		l.Info("[Startup TLS] server key file given", zap.String("key_file", f.ServerKeyFile))
		keyBytes, err = os.ReadFile(f.ServerKeyFile)
		if err != nil {
			return fmt.Errorf("failed to read key file %s: %v", f.ServerKeyFile, err)
		}
	} else {
		return fmt.Errorf("no server key gven for tls")
	}

	if len(keyBytes) > 0 && len(certBytes) > 0 {
		l.Info("[Startup TLS] Loading key pair...")
		tlsCert, err = tls.X509KeyPair(certBytes, keyBytes)
		if err != nil {
			return fmt.Errorf("failed to load key pair: %v", err)
		}
		cfg.Certificates = append(cfg.Certificates, tlsCert)
	}

	ctx.tlsConfig = cfg

	return nil
}

func (a *app) _startup_rsa(ctx *StartupContext) error {
	l := ctx.L()

	privKeyPath := a.features.RSA.PrivateKeyPath
	pubKeyPath := a.features.RSA.PublicKeyPath

	if privKeyPath != "" && pubKeyPath != "" {
		l.Debug("[Startup RSA] initializing RSA keys with paths",
			zap.String("rsa_private_key_path", privKeyPath),
			zap.String("rsa_pub_key_path", pubKeyPath),
		)
		if !fileExists(privKeyPath) {
			l.Info("[Startup RSA] RSA private key not found", zap.String("path", privKeyPath))
			return nil
		}

		if !fileExists(pubKeyPath) {
			l.Info("[Startup RSA] RSA public key not found", zap.String("path", pubKeyPath))
			return nil
		}
		privKey := mustReadFile(privKeyPath)
		pubKey := mustReadFile(pubKeyPath)

		if err := keys.InitRSA(privKey, pubKey); err != nil {
			return err
		}

		a.state.RSAInitialized = true
		l.Debug("[Startup RSA] RSA keys initialized successfully")
	} else {
		l.Debug("[Startup RSA] no RSA key paths provided")
		rsa, err := keys.NewRSA()
		if err != nil {
			return err
		}
		keys.SetRSA(rsa)
	}

	return nil
}

func (a *app) _startup_jwt(ctx *StartupContext) error {
	l := ctx.L()

	configPaths := a.features.JWT.tokenConfigurationPaths
	privKeyPath := a.features.JWT.PrivateKeyPath
	pubKeyPath := a.features.JWT.PubKeyPath

	if privKeyPath != "" && pubKeyPath != "" {

		if !fileExists(privKeyPath) {
			l.Info("[Startup JWT] JWT private key does not exist", zap.String("path", privKeyPath))
			return nil
		}

		if !fileExists(pubKeyPath) {
			l.Info("[Startup JWT] JWT public key does not exist", zap.String("path", pubKeyPath))
			return nil
		}
		l.Debug("[Startup JWT] initializing JWT keys with paths",
			zap.String("jwt_private_key_path", privKeyPath),
			zap.String("jwt_pub_key_path", pubKeyPath),
		)
		jwtPrivKey := mustReadFile(privKeyPath)
		jwtPubKey := mustReadFile(pubKeyPath)
		if err := keys.InitJwt(jwtPrivKey, jwtPubKey); err != nil {
			return err
		}

		a.state.JWTInitialized = true
		l.Debug("[Startup JWT] JWT keys initialized successfully")
	} else {
		l.Debug("[Startup JWT] no JWT key paths provided, creating a key...")
		rsaKey, err := keys.NewRSA()
		if err != nil {
			return err
		}
		jwtKey := keys.NewJWTKey(*rsaKey)
		keys.SetJwt(jwtKey)
	}

	for _, configPath := range configPaths {
		if !fileExists(configPath) {
			l.Error("[Startup JWT] JWT token config file not found", zap.String("path", configPath))
			return fmt.Errorf("JWT token config file not found: %s", configPath)
		}

		config, err := jwt.ParseTokenConfigFile(configPath)
		if err != nil {
			return fmt.Errorf("failed to parse token config file %s: %v", configPath, err)
		}
		ctx.issuerToTokenConfigs[config.Issuer] = *config
	}

	return nil
}

func (a *app) _startup_registry(ctx *StartupContext) error {
	l := ctx.L()

	if a.features.Registry.registryPath != nil && *a.features.Registry.registryPath != "" {
		registryPath := *a.features.Registry.registryPath
		if !fileExists(registryPath) {
			l.Info("[Startup Registry] registry file not found", zap.String("registry_path", registryPath))
			return ErrRegistryFileNotFound
		}

		l.Debug("[Startup Registry] initializing registry with path", zap.String("path", registryPath))
		if err := registry.Init(registryPath); err != nil {
			l.Error("[Startup Registry] failed to initialize registry", zap.Error(err))
			return err
		}
		a.state.RegistryInitialized = true
		l.Debug("[Startup Registry] registry initialized successfully")

	} else {
		l.Debug("[Startup Registry] no registry path provided, using localhost")
		registry.InitLocalhost()
	}

	return nil

}

func (a *app) _startup() error {
	l := a.l

	if a.onPanic != nil {
		defer func() {
			if err := recover(); err != nil {
				l.Warn("recovered from panic", zap.Any("error", err))
				a.onPanic(err)
			}
		}()
	}

	if a.preStartup != nil {
		l.Debug("[Startup] running pre-startup function")
		a.preStartup()
		l.Debug("[Startup] pre-startup function completed")
	}

	startup_funcs := []func(ctx *StartupContext) error{}

	if a.features.Registry.enabled {
		l.Info("[Startup] Enabled registry")
		startup_funcs = append(startup_funcs, a._startup_registry)
	}

	if a.features.JWT.Enabled {
		l.Info("[Startup] JWT enabled")
		startup_funcs = append(startup_funcs, a._startup_jwt)
	}

	if a.features.RSA.Enabled {
		l.Info("[Startup] RSA enabled")
		startup_funcs = append(startup_funcs, a._startup_rsa)
	}

	if a.features.SQL.Enabled {
		l.Info("[Startup] SQL enabled")
		startup_funcs = append(startup_funcs, a._startup_sql)
	}

	if a.features.HTTP.Enabled {
		l.Info("[Startup] HTTP enabled")
		startup_funcs = append(startup_funcs, a._startup_http)
	}

	if a.features.Docs.Enabled {
		l.Info("[Startup] Docs enabled")
		startup_funcs = append(startup_funcs, a._startup_docs)
	}

	if a.features.TLS.Enabled {
		l.Info("[Startup] Gin enabled")
		startup_funcs = append(startup_funcs, a._startup_tls)
	}

	ctx := NewStartupContext(context.Background(), a.l, a.e)
	for _, f := range startup_funcs {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	l.Debug("[Startup] adding logging routes")
	v1.AddRoutes(a.e)
	l.Debug("[Startup] finished adding logging routes")
	if a.startup != nil {
		l.Debug("[Startup] Running app...")
		err := a.startup(ctx)
		if err != nil {
			l.Error("[Startup] encountered an error on startup", zap.Error(err))
			return err
		}
		l.Debug("[Startup] app exited")

	}

	return nil
}
