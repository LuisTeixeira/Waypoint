package memory

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/luisteixeira/waypoint/backend/internal/domain"
	"github.com/luisteixeira/waypoint/backend/internal/repository"
)

type InMemoryDefintionRepo struct {
	mu          sync.RWMutex
	definitions map[uuid.UUID]domain.ActivityDefinition
}

func NewInMemoryDefinitionRepo() *InMemoryDefintionRepo {
	return &InMemoryDefintionRepo{
		definitions: make(map[uuid.UUID]domain.ActivityDefinition),
	}
}

func (r *InMemoryDefintionRepo) GetOrCreateByName(ctx context.Context, name string) (*domain.ActivityDefinition, error) {
	familyID, err := repository.GetFamilyIdFromContext(ctx)
	if err != nil {
		return nil, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, d := range r.definitions {
		if d.FamilyID == familyID && d.Name == name {
			return &d, nil
		}
	}

	newDef := domain.ActivityDefinition{
		ID:       uuid.New(),
		FamilyID: familyID,
		Name:     name,
	}
	r.definitions[newDef.ID] = newDef
	return &newDef, nil
}

func (r *InMemoryDefintionRepo) ListByFamily(ctx context.Context) ([]domain.ActivityDefinition, error) {
	familyID, err := repository.GetFamilyIdFromContext(ctx)
	if err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.Unlock()

	var defs []domain.ActivityDefinition
	for _, d := range r.definitions {
		if d.FamilyID == familyID {
			defs = append(defs, d)
		}
	}

	return defs, nil
}
