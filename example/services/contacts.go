package services

import (
	"context"

	"github.com/bilal-bhatti/zipline/example/connectors"
	"github.com/bilal-bhatti/zipline/example/models"
)

type ContactsService struct {
	DataConnector *connectors.DataConnector
}

func NewContactsService(dc *connectors.DataConnector) ContactsService {
	return ContactsService{
		DataConnector: dc,
	}
}

func (cs ContactsService) Create(ctx context.Context, contactRequest *models.ContactRequest) (*models.ContactResponse, error) {
	return &models.ContactResponse{Output: "Out"}, nil
}

func (cs ContactsService) Update(ctx context.Context, id int, contactRequest models.ContactRequest) (*models.ContactResponse, error) {
	return &models.ContactResponse{Output: "Out"}, nil
}

func (cs ContactsService) Replace(ctx context.Context, id int, contactRequest models.ContactRequest) (*models.ContactResponse, error) {
	return &models.ContactResponse{Output: "Out"}, nil
}

func (cs ContactsService) GetOne(ctx context.Context, id int) (*models.ContactResponse, error) {
	return &models.ContactResponse{Output: "Out"}, nil
}

func (cs ContactsService) GetByDate(ctx context.Context, month, day, year string) (*models.ContactResponse, error) {
	return &models.ContactResponse{Output: "Out"}, nil
}
