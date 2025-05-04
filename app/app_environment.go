package app

import (
	"context"
	"time"

	"github.com/ooqls/go-db/testutils"
	"github.com/ooqls/go-log"
	"github.com/testcontainers/testcontainers-go"
	"go.uber.org/zap"
)

type TestEnvironment struct {
	Postgres bool
	Redis    bool
}

func (e *TestEnvironment) Start(ctx context.Context) (func(), error) {
	l := log.NewLogger("testEnvironment")
	var containers []testcontainers.Container
	if e.Postgres {
		cont := testutils.StartPostgres(ctx)
		containers = append(containers, cont)
	}

	if e.Redis {
		redisCont := testutils.InitRedis()
		containers = append(containers, redisCont)
	}

	return func() {
		for _, c := range containers {
			timeout := time.Second * 30
			err := c.Stop(ctx, &timeout)
			if err != nil {
				l.Error("failed to stop container", zap.Error(err))
			}
		}
	}, nil
}
