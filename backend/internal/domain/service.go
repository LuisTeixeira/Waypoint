package domain

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrEntityBusy = errors.New("child is already participating in an activity")

type StartActivityInput struct {
	EntityID           uuid.UUID
	DefinitionID       uuid.UUID
	NewDefinittionName string
	CaregiversIDs      []uuid.UUID
}

type ActivityService interface {
	StartActivity(ctx context.Context, input StartActivityInput) (ActivityRealization, error)
	CompleteActivity(ctx context.Context, realizationID uuid.UUID) error
	PlanActivity(ctx context.Context, input StartActivityInput) (*ActivityRealization, error)
}
