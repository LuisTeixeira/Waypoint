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
	defRepo := memory.NewInMemoryDefinitionRepo()
	svc := service.NewActivityService(repo, defRepo)

	familyID := uuid.New()
	entityID := uuid.New()
	ctx := context.WithValue(context.Background(), middleware.FamilyIDKey, familyID)

	t.Run("Successfully start activity when child is free", func(t *testing.T) {
		defID := uuid.New()
		caregivers := []uuid.UUID{uuid.New()}

		input := domain.StartActivityInput{
			DefinitionID:  defID,
			EntityID:      entityID,
			CaregiversIDs: caregivers,
		}

		ar, err := svc.StartActivity(ctx, input)

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

		input := domain.StartActivityInput{
			DefinitionID:  defID,
			EntityID:      entityID,
			CaregiversIDs: caregivers,
		}

		ar, err := svc.StartActivity(ctx, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrEntityBusy)
		assert.Nil(t, ar)
	})
}

func TestActivityService_PlanActivity(t *testing.T) {
	repo := memory.NewInMemoryActivityRepo()
	defRepo := memory.NewInMemoryDefinitionRepo()
	svc := service.NewActivityService(repo, defRepo)

	familyID := uuid.New()
	entityID := uuid.New()
	ctx := context.WithValue(context.Background(), middleware.FamilyIDKey, familyID)

	t.Run("Plan activity with new definition name", func(t *testing.T) {
		input := domain.StartActivityInput{
			EntityID:           entityID,
			NewDefinittionName: "Sport",
			CaregiversIDs:      []uuid.UUID{uuid.New()},
		}

		activityRealization, err := svc.PlanActivity(ctx, input)

		assert.NoError(t, err)
		assert.NotNil(t, activityRealization)
		assert.Equal(t, domain.StatusPlanned, activityRealization.Status)
		assert.Nil(t, activityRealization.StartedAt, "Planned activities should not have a started at")

		defs, _ := defRepo.ListByFamily(ctx)
		assert.Len(t, defs, 1)
		assert.Equal(t, "Sport", defs[0].Name)
	})

	t.Run("Plan activity even if child is currently busy", func(t *testing.T) {
		_, _ = svc.StartActivity(ctx, domain.StartActivityInput{
			EntityID:           entityID,
			NewDefinittionName: "Dinner",
		})

		input := domain.StartActivityInput{
			EntityID:           entityID,
			NewDefinittionName: "Sleep",
		}

		activityRealization, err := svc.PlanActivity(ctx, input)

		assert.NoError(t, err, "Planning should not be blocked by activity in progress")
		assert.Equal(t, domain.StatusPlanned, activityRealization.Status)
	})

	t.Run("Fail if neither ID nor Name is provided", func(t *testing.T) {
		input := domain.StartActivityInput{
			EntityID: entityID,
		}

		activityRealization, err := svc.PlanActivity(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, activityRealization)
		assert.Contains(t, err.Error(), "either definition_id or new_definition_name")
	})
}
