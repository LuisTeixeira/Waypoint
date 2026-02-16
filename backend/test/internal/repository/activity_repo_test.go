package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/luisteixeira/waypoint/backend/internal/domain"
	"github.com/luisteixeira/waypoint/backend/internal/middleware"
	"github.com/luisteixeira/waypoint/backend/test/internal/repository/memory"
	"github.com/stretchr/testify/assert"
)

func TestActivityRepository(t *testing.T) {
	repo := memory.NewInMemoryActivityRepo()
	familyID := uuid.New()

	ctx := context.WithValue(context.Background(), middleware.FamilyIDKey, familyID)

	entityID := uuid.New()
	caregiverIDs := []uuid.UUID{uuid.New(), uuid.New()}
	startTime := time.Now().Round(time.Second)

	input := &domain.ActivityRealization{
		DefinitionID:  uuid.New(),
		EntityID:      entityID,
		CaregiversIDs: caregiverIDs,
		Status:        "in_progress",
		StartedAt:     &startTime,
	}

	t.Run("Create and Retrieve Realization", func(t *testing.T) {
		err := repo.CreateRealization(ctx, input)
		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, input.ID, "Repo should assing a new id")

		output, err := repo.GetRealizationByID(ctx, input.ID)

		assert.NoError(t, err)
		assert.Equal(t, input.ID, output.ID)
		assert.Equal(t, familyID, output.FamilyID)
		assert.Equal(t, "in_progress", output.Status)

		assert.Len(t, output.CaregiversIDs, 2)
		assert.Contains(t, output.CaregiversIDs, caregiverIDs[0])
		assert.Contains(t, output.CaregiversIDs, caregiverIDs[1])
	})

	t.Run("Multi-tenancy Protection", func(t *testing.T) {
		wrongFamilyCtx := context.WithValue(context.Background(), middleware.FamilyIDKey, uuid.New())

		err := repo.CreateRealization(ctx, input)
		assert.NoError(t, err)
		output, err := repo.GetRealizationByID(wrongFamilyCtx, input.ID)

		assert.Error(t, err, "Should fail when accessing data from another family")
		assert.Nil(t, output)
	})
}
