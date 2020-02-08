package services

import (
	"context"
	"net/http"

	"github.com/bilal-bhatti/zipline/example/connectors"
	"github.com/bilal-bhatti/zipline/example/models"
	"github.com/go-chi/chi"
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

func (cs ContactsService) Replace(ctx context.Context, id int, contactRequest models.ContactRequest) (*models.ContactResponse, error) {
	return &models.ContactResponse{Output: "Out"}, nil
}

func (cs ContactsService) GetOne(ctx context.Context, id int) (*models.ContactResponse, error) {
	return &models.ContactResponse{Output: "Out"}, nil
}

type DateFilter struct {
	Month, Day, Year string
}

func ProvideDateFilter(request *http.Request) DateFilter {
	return DateFilter{
		Month: chi.URLParam(request, "month"),
		Day:   chi.URLParam(request, "day"),
		Year:  chi.URLParam(request, "year"),
	}
}

func (cs ContactsService) GetByDate(df DateFilter) (*models.ContactResponse, error) {
	return &models.ContactResponse{Output: "Out"}, nil
}
