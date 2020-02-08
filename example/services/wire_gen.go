// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package services

import (
	"github.com/bilal-bhatti/zipline/example/connectors"
)

// Injectors from wire.go:

func InitContactsService() *ContactsService {
	dataConnector := connectors.ProvideDataConnector()
	contactsService := &ContactsService{
		DataConnector: dataConnector,
	}
	return contactsService
}

func InitThingsService() *ThingsService {
	dataConnector := connectors.ProvideDataConnector()
	thingsService := &ThingsService{
		DataConnector: dataConnector,
	}
	return thingsService
}
