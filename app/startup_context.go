package app

import (
	"context"

	"go.uber.org/zap"
)

func NewStartupContext(ctx context.Context, l *zap.Logger) *StartupContext {
	return &StartupContext{
		l:       l,
		Context: ctx,
	}
}

type StartupContext struct {
	context.Context
	l *zap.Logger
}

func (ctx *StartupContext) L() *zap.Logger {
	return ctx.l
}
