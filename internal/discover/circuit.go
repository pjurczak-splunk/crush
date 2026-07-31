package discover

import (
	"context"

	"charm.land/catwalk/pkg/catwalk"
)

func init() {
	RegisterEnricher("circuit", &circuitEnricher{})
}

// circuitEnricher leaves the model list unchanged.
//
// Circuit is validated as a first-class custom provider type, but the
// deployment IDs still need to be declared explicitly in config. We do
// not currently have a richer discovery endpoint to call here.
type circuitEnricher struct{}

func (e *circuitEnricher) EnrichModels(_ context.Context, _ Config, _ Resolver, models []catwalk.Model) ([]catwalk.Model, error) {
	return models, nil
}
