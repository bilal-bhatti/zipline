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

// Create thing
func (cs ThingsService) Create(ctx context.Context, req models.ThingRequest) (*models.ThingResponse, error) {
	return &models.ThingResponse{}, nil
}

// Delete thing by id
func (cs *ThingsService) Delete(id int) error {
	return nil
}

// Get things by category and search query
//
// @parameters  {"name":"category", "description": "category of data to search"}
// @parameters  {"name":"q", "description": "search query"}
func (cs ThingsService) GetByCategoryAndQuery(ctx context.Context, category string, q string) (ThingListResponse, error) {
	log.Println("Getting by category and query", category, q)
	return ThingListResponse{}, nil
}

// Get things by date range
//
// @description       A long description of this endpoint
// @summary           A short summary of this endpoint
// @operationId       GetThingsByDateRange
// @consumes          application/json
// @produces          application/json
// @produces          plain/text
// @tags              ["things", "example", "get"]
// @parameters        {"name":"from", "format":"date-time,2006-01-02", "description": "date should be in Go time format"}
// @parameters        {"name":"to", "format":"date-time,2006-01-02", "description": "date should be in Go time format"}
// @parameters        {"name":"notgood", "description": "parameter not found in code, tsk tsk", "in": "path", "type": "string", "format": "eMail"}
// @responses.400     {models.ErrorResponse}
// @responses.404     {models.ErrorResponse}
func (cs ThingsService) GetByDateRange(ctx context.Context, from, to *time.Time) (ThingListResponse, error) {
	log.Println("Getting by category and query", from, to)
	return ThingListResponse{}, nil
}
