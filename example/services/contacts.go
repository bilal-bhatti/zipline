package services

import (
	"context"

	"github.com/bilal-bhatti/zipline/example/connectors"
)

type (
	Name struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}
	Address struct {
		Street  string `json:"street"`
		City    string `json:"city"`
		State   string `json:"state"`
		ZipCode string `json:"zipCode"`
	}

	// ContactRequest model
	ContactRequest struct {
		ID string `json:"id"`
		Name
		EMail   string   `json:"eMail" format:"email"`
		Address *Address `json:"address"`
	}

	// ContactResponse model
	ContactResponse struct {
		// id is the unique id of contact
		ID string `json:"id"`
	}

	ContactsService struct {
		DataConnector *connectors.DataConnector
	}
)

func NewContactsService(dc *connectors.DataConnector) ContactsService {
	return ContactsService{
		DataConnector: dc,
	}
}

// Create a new contact request entity.
func (cs ContactsService) Create(ctx context.Context, contactRequest *ContactRequest) (*ContactResponse, error) {
	return &ContactResponse{ID: "id"}, nil
}

// Update a contact entity with provided data.
func (cs ContactsService) Update(ctx context.Context, id int, contactRequest ContactRequest) (*ContactResponse, error) {
	return &ContactResponse{ID: "id"}, nil
}

// Replace a contact entity completely.
func (cs *ContactsService) Replace(ctx context.Context, id int, contactRequest ContactRequest) (*ContactResponse, error) {
	return &ContactResponse{ID: "id"}, nil
}

// GetOne contact by id
// id contact id
func (cs ContactsService) GetOne(ctx context.Context, id int) (*ContactResponse, error) {
	return &ContactResponse{ID: "id"}, nil
}

// GetBunch of contacts by ids
//
// @summary          Get a list of contacts by ids
// @description      Get a list of contacts by ids
// @tags             contacts
// @produces         application/json
// @parameters       {"name":"ids", "description":"list of contact ids", "required":true}
// @responses.400    {models.ErrorResponse}
func (cs ContactsService) GetBunch(ctx context.Context, ids []int64) (*ContactResponse, error) {
	return &ContactResponse{ID: "id"}, nil
}

// Get contacts list by date
func (cs ContactsService) GetByDate(ctx context.Context, month, day, year string) (*ContactResponse, error) {
	return &ContactResponse{ID: "id"}, nil
}

// DeleteBulk contact by id
func (cs ContactsService) DeleteBulk(ctx context.Context, ids []string) error {
	return nil
}

// Redirect sends a redirect in response
func (cs ContactsService) Redirect(ctx context.Context, id string) (string, error) {
	// return url.Parse("https://donate.doctorswithoutborders.org")
	return "https://donate.doctorswithoutborders.org", nil
}
