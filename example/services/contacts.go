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

// Create a new contact request entity.
func (cs ContactsService) Create(ctx context.Context, contactRequest *models.ContactRequest) (*models.ContactResponse, error) {
	return &models.ContactResponse{ID: "id"}, nil
}

// Update a contact entity with provided data.
func (cs ContactsService) Update(ctx context.Context, id int, contactRequest models.ContactRequest) (*models.ContactResponse, error) {
	return &models.ContactResponse{ID: "id"}, nil
}

// Replace a contact entity completely.
func (cs *ContactsService) Replace(ctx context.Context, id int, contactRequest models.ContactRequest) (*models.ContactResponse, error) {
	return &models.ContactResponse{ID: "id"}, nil
}

// GetOne contact by id
func (cs ContactsService) GetOne(ctx context.Context, id int) (*models.ContactResponse, error) {
	return &models.ContactResponse{ID: "id"}, nil
}

// Get contacts list by date
func (cs ContactsService) GetByDate(ctx context.Context, month, day, year string) (*models.ContactResponse, error) {
	return &models.ContactResponse{ID: "id"}, nil
}

// DeleteBulk contact by id
func (cs ContactsService) DeleteBulk(ctx context.Context, ids []string) error {
	return nil
}
