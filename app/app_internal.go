package app

import (
	"context"
	"fmt"

	"github.com/ooqls/go-crypto/keys"
	v1 "github.com/ooqls/go-log/api/v1"
	"github.com/ooqls/go-registry"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (a *app) _startup_rsa(ctx *StartupContext) error {
	l := ctx.L().WithOptions(zap.Hooks(func(e zapcore.Entry) error {
		e.Message = fmt.Sprintf("[Startup RSA] %s", e.Message)
		return nil
	}))
	
	privKeyPath := a.features.RSA.PrivateKeyPath
	pubKeyPath := a.features.RSA.PublicKeyPath

	if privKeyPath != "" && pubKeyPath != "" {
		l.Debug("initializing RSA keys with paths",
			zap.String("rsa_private_key_path", privKeyPath),
			zap.String("rsa_pub_key_path", pubKeyPath),
		)
		if !fileExists(privKeyPath) {
			l.Info("RSA private key not found", zap.String("path", privKeyPath))
			return nil
		}

		if !fileExists(pubKeyPath) {
			l.Info("RSA public key not found", zap.String("path", pubKeyPath))
			return nil
		}
		privKey := mustReadFile(privKeyPath)
		pubKey := mustReadFile(pubKeyPath)

		if err := keys.InitRSA(privKey, pubKey); err != nil {
			return err
		}

		a.state.RSAInitialized = true
		l.Debug("RSA keys initialized successfully")
	} else {
		l.Debug("no RSA key paths provided")
	}

	return nil
}

func (a *app) _startup_jwt(ctx *StartupContext) error {
	l := ctx.L().WithOptions(zap.Hooks(func(e zapcore.Entry) error {
		e.Message = fmt.Sprintf("[Startup JWT] %s", e.Message)
		return nil
	}))

	privKeyPath := a.features.JWT.PrivateKeyPath
	pubKeyPath := a.features.JWT.PubKeyPath

	if privKeyPath != "" && pubKeyPath != "" {

		if !fileExists(privKeyPath) {
			l.Info("JWT private key does not exist", zap.String("path", privKeyPath))
			return nil
		}

		if !fileExists(pubKeyPath) {
			l.Info("JWT public key does not exist", zap.String("path", pubKeyPath))
			return nil
		}
		l.Debug("initializing JWT keys with paths",
			zap.String("jwt_private_key_path", privKeyPath),
			zap.String("jwt_pub_key_path", pubKeyPath),
		)
		jwtPrivKey := mustReadFile(privKeyPath)
		jwtPubKey := mustReadFile(pubKeyPath)
		if err := keys.InitJwt(jwtPrivKey, jwtPubKey); err != nil {
			return err
		}

		a.state.JWTInitialized = true
		l.Debug("JWT keys initialized successfully")
	} else {
		l.Debug("no JWT key paths provided")
	}

	return nil
}

func (a *app) _startup_registry(ctx *StartupContext) error {
	l := ctx.L().WithOptions(zap.Hooks(func(e zapcore.Entry) error {
		e.Message = fmt.Sprintf("[Startup Registry] %s", e.Message)
		return nil
	}))
	registryPath := a.features.Registry.registryPath
	if registryPath != "" {
		if !fileExists(registryPath) {
			l.Info("registry file not found", zap.String("registry_path", registryPath))
			return nil
		}

		l.Debug("initializing registry with path", zap.String("path", registryPath))
		if err := registry.Init(registryPath); err != nil {
			l.Error("failed to initialize registry", zap.Error(err))
			return err
		}
		a.state.RegistryInitialized = true
		l.Debug("registry initialized successfully")

	} else {
		l.Debug("no registry path provided")
	}

	return nil

}

func (a *app) _startup() error {
	l := a.l.WithOptions(zap.Hooks(func(e zapcore.Entry) error {
		e.Message = fmt.Sprintf("[Startup] %s", e.Message)
		return nil
	}))

	if a.onPanic != nil {
		defer func() {
			if err := recover(); err != nil {
				l.Warn("recovered from panic", zap.Any("error", err))
				a.onPanic(err)
			}
		}()
	}

	if a.preStartup != nil {
		l.Debug("running pre-startup function")
		a.preStartup()
		l.Debug("pre-startup function completed")
	}

	startup_funcs := []func(ctx *StartupContext) error{
		a._startup_registry,
	}

	if a.features.JWT.Enabled {
		startup_funcs = append(startup_funcs, a._startup_jwt)
	}

	if a.features.RSA.Enabled {
		startup_funcs = append(startup_funcs, a._startup_rsa)
	}

	if a.features.SQL.Enabled {
		startup_funcs = append(startup_funcs, a._startup_sql)
	}

	ctx := NewStartupContext(context.Background(), a.l)
	for _, f := range startup_funcs {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	l.Debug("adding logging routes")
	v1.AddRoutes(a.e)
	l.Debug("finished adding logging routes")

	l.Info("App state", zap.Dict("state",
		zap.Bool("RSAInitialized", a.state.RSAInitialized),
		zap.Bool("JWTInitialized", a.state.JWTInitialized),
		zap.Bool("SQLInitialized", a.state.SQLInitialized),
		zap.Bool("SQLSeeded", a.state.SQLSeeded),
		zap.Bool("RegistryInitialized", a.state.RegistryInitialized)))
	if a.startup != nil {
		l.Debug("running startup function")
		err := a.startup(a.e)
		if err != nil {
			l.Error("encountered an error on startup", zap.Error(err))
			return err
		}
		l.Debug("startup function completed")

	}

	return nil
}
