package web

import (
	"context"
	"log"
	"time"

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

// Get things by date range
//
// from `format:"date-time,2006-01-02"` date should be in Go time format"
// @to `format:"date-time,2006-01-02"` date should be in Go time format"
func (cs ThingsService) GetByDateRange(ctx context.Context, from, to time.Time) (ThingListResponse, error) {
	log.Println("Getting by category and query", from, to)
	return ThingListResponse{}, nil
}
