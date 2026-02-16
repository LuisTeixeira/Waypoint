package domain

import (
	"context"

	"github.com/google/uuid"
)

type ActivityRepository interface {
	CreateRealization(ctx context.Context, activityRealization *ActivityRealization) error
	GetRealizationByID(ctx context.Context, id uuid.UUID) (*ActivityRealization, error)
}
