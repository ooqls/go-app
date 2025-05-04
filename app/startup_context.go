package app

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ooqls/go-crypto/jwt"
	"go.uber.org/zap"
)

func NewStartupContext(ctx context.Context, l *zap.Logger, e *gin.Engine) *StartupContext {
	return &StartupContext{
		l:       l,
		Context: ctx,
		e:       e,
	}
}

type StartupContext struct {
	context.Context
	l          *zap.Logger
	httpClient *http.Client
	tlsConfig  *tls.Config
	e          *gin.Engine
	tokenConfig *jwt.TokenConfiguration
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

func (ctx *StartupContext) IssuerConfig() (*jwt.TokenConfiguration, bool) {
	return ctx.tokenConfig, ctx.tokenConfig != nil
}
