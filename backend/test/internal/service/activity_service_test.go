package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/luisteixeira/waypoint/backend/internal/domain"
	"github.com/luisteixeira/waypoint/backend/internal/middleware"
	"github.com/luisteixeira/waypoint/backend/internal/service"
	"github.com/luisteixeira/waypoint/backend/test/internal/repository/memory"
	"github.com/stretchr/testify/assert"
)

func TestActivityService_StartActivity(t *testing.T) {
	repo := memory.NewInMemoryActivityRepo()
	svc := service.NewActivityService(repo)

	familyID := uuid.New()
	entityID := uuid.New()
	ctx := context.WithValue(context.Background(), middleware.FamilyIDKey, familyID)

	t.Run("Successfully start activity when child is free", func(t *testing.T) {
		defID := uuid.New()
		caregivers := []uuid.UUID{uuid.New()}

		ar, err := svc.StartActivity(ctx, defID, entityID, caregivers)

		assert.NoError(t, err)
		assert.NotNil(t, ar)
		assert.Equal(t, domain.StatusInProgress, ar.Status)
		assert.NotNil(t, ar.StartedAt)
	})

	t.Run("Fail to start activity when child is already busy", func(t *testing.T) {
		// Child is already in the activity from the previous test case
		// (since we are reusing the same 'repo' and 'entityID')

		defID := uuid.New()
		caregivers := []uuid.UUID{uuid.New()}

		ar, err := svc.StartActivity(ctx, defID, entityID, caregivers)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrEntityBusy)
		assert.Nil(t, ar)
	})
}
