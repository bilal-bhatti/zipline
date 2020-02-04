package services

import (
	"context"

	"github.com/bilal-bhatti/zipline/example/connectors"
	"github.com/bilal-bhatti/zipline/example/models"
)

type ContactsService struct {
	DataConnector *connectors.DataConnector
}

func (cs ContactsService) Create(ctx context.Context, contactRequest *models.ContactRequest) (*models.ContactResponse, error) {
	return &models.ContactResponse{Output: "Out"}, nil
}

func (cs ContactsService) Update(ctx context.Context, id int, contactRequest models.ContactRequest) (*models.ContactResponse, error) {
	return &models.ContactResponse{Output: "Out"}, nil
}

func (cs ContactsService) GetOne(id int) (*models.ContactResponse, error) {
	return &models.ContactResponse{Output: "Out"}, nil
}
