// +build wireinject

package services

import (
	"github.com/bilal-bhatti/zipline/example/connectors"
	"github.com/google/wire"
)

func InitContactsService() *ContactsService {
	panic(wire.Build(connectors.ProvideDataConnector, wire.Struct(new(ContactsService), "*")))
}

func InitThingsService() *ThingsService {
	panic(wire.Build(connectors.ProvideDataConnector, wire.Struct(new(ThingsService), "*")))
}
