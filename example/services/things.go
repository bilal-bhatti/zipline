package services

import (
	"context"

	"github.com/bilal-bhatti/zipline/example/connectors"
)

type ThingsService struct {
	DataConnector *connectors.DataConnector
}

func (cs ThingsService) Delete(ctx context.Context, id int) error {
	return nil
}
