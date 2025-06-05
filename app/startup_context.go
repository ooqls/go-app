package app

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ooqls/go-crypto/jwt"
	"go.uber.org/zap"
)

const (
	AuthIssuer    = "auth"
	RefreshIssuer = "refresh"
)

func NewStartupContext(ctx context.Context, l *zap.Logger, e *gin.Engine) *StartupContext {
	return &StartupContext{
		l:                    l,
		Context:              ctx,
		e:                    e,
		issuerToTokenConfigs: make(map[string]jwt.TokenConfiguration),
	}
}

type StartupContext struct {
	context.Context
	l                    *zap.Logger
	httpClient           *http.Client
	tlsConfig            *tls.Config
	e                    *gin.Engine
	issuerToTokenConfigs map[string]jwt.TokenConfiguration
}

func (ctx *StartupContext) L() *zap.Logger {
	return ctx.l
}

func (ctx *StartupContext) HTTPClient() http.Client {
	if ctx.httpClient == nil {
		return *http.DefaultClient
	}

	return *ctx.httpClient
}

func (ctx *StartupContext) TLSConfig() (*tls.Config, bool) {
	return ctx.tlsConfig, ctx.tlsConfig != nil
}

func (ctx *StartupContext) Gin() (*gin.Engine, bool) {
	return ctx.e, ctx.e != nil
}

func (ctx *StartupContext) AuthIssuerConfig() (*jwt.TokenConfiguration, bool) {
	config, ok := ctx.issuerToTokenConfigs[AuthIssuer]
	return &config, ok
}

func (ctx *StartupContext) RefreshIssuerConfig() (*jwt.TokenConfiguration, bool) {
	config, ok := ctx.issuerToTokenConfigs[RefreshIssuer]
	return &config, ok
}
