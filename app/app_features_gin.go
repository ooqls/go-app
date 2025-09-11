package app

import "github.com/gin-gonic/gin"

const (
	gin_engineOpt string = "opt-gin-engine"
	gin_portOpt   string = "opt-gin-port"
)

type ginOpt struct {
	featureOpt
}

func WithGin(e *gin.Engine) ginOpt {
	return ginOpt{
		featureOpt: featureOpt{
			key:   gin_engineOpt,
			value: e,
		},
	}
}

type GinFeature struct {
	Enabled bool
	Port    int
	Engine  *gin.Engine
}

func (f *GinFeature) apply(opt ginOpt) {
	switch opt.key {
	case gin_engineOpt:
		f.Engine = opt.value.(*gin.Engine)
	case gin_portOpt:
		f.Port = opt.value.(int)
	}
}

func Gin(opts ...ginOpt) GinFeature {
	f := GinFeature{
		Enabled: true,
		Port:    8080,
		Engine:  gin.New(),
	}

	for _, opt := range opts {
		f.apply(opt)
	}

	return f
}
