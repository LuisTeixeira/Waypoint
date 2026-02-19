package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/luisteixeira/waypoint/backend/internal/middleware"
)

func GetFamilyIdFromContext(ctx context.Context) (uuid.UUID, error) {
	familyID, ok := ctx.Value(middleware.FamilyIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("unauthorized: family_id missing")
	}
	return familyID, nil
}
