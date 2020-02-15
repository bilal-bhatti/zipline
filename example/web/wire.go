// +build wireinject

package web

import (
	"github.com/bilal-bhatti/zipline/example/connectors"
	"github.com/google/wire"
)

func InitThingsService() (*ThingsService, error) {
	panic(wire.Build(connectors.ProvideDataConnector, wire.Struct(new(ThingsService), "*")))
}
