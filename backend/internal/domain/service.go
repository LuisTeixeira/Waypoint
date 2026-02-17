package domain

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrEntityBusy = errors.New("child is already participating in an activity")

type ActivityService interface {
	StartActivity(ctx context.Context, definitionID, entityID uuid.UUID, caregiverIDs []uuid.UUID) (ActivityRealization, error)
	CompleteActivity(ctx context.Context, realizationID uuid.UUID) error
}
