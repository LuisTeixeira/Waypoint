package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/luisteixeira/waypoint/backend/internal/domain"
	"github.com/luisteixeira/waypoint/backend/internal/middleware"
)

type InMemoryActivityRepo struct {
	mu           sync.RWMutex
	realizations map[uuid.UUID]domain.ActivityRealization
}

func NewInMemoryActivityRepo() *InMemoryActivityRepo {
	return &InMemoryActivityRepo{
		realizations: make(map[uuid.UUID]domain.ActivityRealization),
	}
}

func (r *InMemoryActivityRepo) CreateRealization(ctx context.Context, activityRealization *domain.ActivityRealization) error {
	familyID, ok := ctx.Value(middleware.FamilyIDKey).(uuid.UUID)
	if !ok {
		return fmt.Errorf("unauthorized: family_id missing")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	activityRealization.ID = uuid.New()
	activityRealization.FamilyID = familyID
	r.realizations[activityRealization.ID] = *activityRealization
	return nil
}

func (r *InMemoryActivityRepo) GetRealizationByID(ctx context.Context, id uuid.UUID) (*domain.ActivityRealization, error) {
	familyID, ok := ctx.Value(middleware.FamilyIDKey).(uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("unauthorized: family_id missing")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	res, ok := r.realizations[id]
	if !ok {
		return nil, fmt.Errorf("Not found")
	}
	if familyID != res.FamilyID {
		return nil, fmt.Errorf("unauthorized: wrong family_id")
	}
	return &res, nil
}

func (r *InMemoryActivityRepo) GetActiveByEntity(ctx context.Context, entityID uuid.UUID) (*domain.ActivityRealization, error) {
	familyID, ok := ctx.Value(middleware.FamilyIDKey).(uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, ar := range r.realizations {
		if ar.FamilyID == familyID &&
			ar.EntityID == entityID &&
			ar.Status == domain.StatusInProgress { // Explicit status check
			copyAr := ar
			return &copyAr, nil
		}
	}
	return nil, nil
}

func (r *InMemoryActivityRepo) UpdateRealization(ctx context.Context, activityRealization *domain.ActivityRealization) error {
	r.mu.Unlock()
	defer r.mu.Unlock()

	existing, ok := r.realizations[activityRealization.ID]
	if !ok || existing.FamilyID != activityRealization.FamilyID {
		return fmt.Errorf("not found")
	}

	r.realizations[activityRealization.ID] = *activityRealization
	return nil
}
