package services

import (
	"context"

	"github.com/bilal-bhatti/zipline/example/connectors"
	"github.com/bilal-bhatti/zipline/example/models"
)

type DoodadsService struct {
	env *connectors.Env
}

func NewDoodadsService(env *connectors.Env) *DoodadsService {
	return &DoodadsService{
		env: env,
	}
}

// Create a new doodad entity.
func (cs DoodadsService) Create(ctx context.Context, contactRequest *models.ThingRequest) (*models.ThingResponse, error) {
	return &models.ThingResponse{Name: "shiny doodad"}, nil
}
