package web

import (
	"context"
	"log"

	"github.com/bilal-bhatti/zipline/example/connectors"
	"github.com/bilal-bhatti/zipline/example/models"
)

type (
	// List of things
	ThingListResponse struct {
		Things []models.ThingResponse `json:"things"`
	}
	ThingsService struct {
		DataConnector *connectors.DataConnector
	}
)

// Delete thing by id
func (cs *ThingsService) Delete(id int) error {
	return nil
}

// Get things by category and search query
func (cs ThingsService) GetByCategoryAndQuery(ctx context.Context, category string, q string) (ThingListResponse, error) {
	log.Println("Getting by category and query", category, q)
	return ThingListResponse{}, nil
}
