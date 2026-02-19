package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/luisteixeira/waypoint/backend/internal/domain"
)

type ActivityRepository interface {
	CreateRealization(ctx context.Context, activityRealization *domain.ActivityRealization) error
	GetRealizationByID(ctx context.Context, id uuid.UUID) (*domain.ActivityRealization, error)
	GetActiveByEntity(ctx context.Context, entityID uuid.UUID) (*domain.ActivityRealization, error)
	UpdateRealization(ctx context.Context, activityRealization *domain.ActivityRealization) error
}
