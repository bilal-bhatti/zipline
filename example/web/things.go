package web

import (
	"context"
	"log"

	"github.com/bilal-bhatti/zipline/example/connectors"
)

type ThingResponse struct {
	Output string `json:"output"`
}
type ThingsService struct {
	DataConnector *connectors.DataConnector
}

// Delete thing by id
func (cs *ThingsService) Delete(id int) error {
	return nil
}

// Get things by category and search query
func (cs ThingsService) GetByCategoryAndQuery(ctx context.Context, category string, q string) ([]ThingResponse, error) {
	log.Println("Getting by category and query", category, q)
	return []ThingResponse{}, nil
}
