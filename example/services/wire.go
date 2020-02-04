// +build wireinject

package services

import (
	"github.com/bilal-bhatti/zipline/example/connectors"
	"github.com/google/wire"
)

func InitContactsService() *ContactsService {
	panic(wire.Build(connectors.ProvideDataConnector, ContactsService{}))
}
