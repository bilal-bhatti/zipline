package services

import (
	"context"
	"net/http"

	"github.com/bilal-bhatti/zipline/example/connectors"
	"github.com/bilal-bhatti/zipline/example/models"
)

type DoodadsService struct {
	env *connectors.Env
}

func NewDoodadsService(env *connectors.Env) (*DoodadsService, error) {
	return &DoodadsService{
		env: env,
	}, nil
}

// Create a new doodad entity.
func (cs DoodadsService) Create(ctx context.Context, r *http.Request, thing *models.ThingRequest) (*models.ThingResponse, error) {
	return &models.ThingResponse{Name: "shiny doodad"}, nil
}
