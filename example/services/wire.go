// +build wireinject

package services

import (
	"github.com/bilal-bhatti/zipline/example/connectors"
	"github.com/google/wire"
)

func InitContactsService() (*ContactsService, error) {
	panic(wire.Build(connectors.ProvideDataConnector, wire.Struct(new(ContactsService), "*")))
}
