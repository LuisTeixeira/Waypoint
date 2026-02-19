package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/luisteixeira/waypoint/backend/internal/domain"
)

type activityService struct {
	repo ActivityRepository
}

func NewActivityService(repo ActivityRepository) *activityService {
	return &activityService{
		repo: repo,
	}
}

func (s *activityService) StartActivity(ctx context.Context, defID, entityID uuid.UUID, caregiverIDs []uuid.UUID) (*domain.ActivityRealization, error) {
	active, err := s.repo.GetActiveByEntity(ctx, entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to check child status: %w", err)
	}

	if active != nil {
		return nil, domain.ErrEntityBusy
	}

	now := time.Now()
	realization := &domain.ActivityRealization{
		DefinitionID:  defID,
		EntityID:      entityID,
		CaregiversIDs: caregiverIDs,
		Status:        domain.StatusInProgress,
		StartedAt:     &now,
	}

	if err := s.repo.CreateRealization(ctx, realization); err != nil {
		return nil, err
	}

	return realization, nil
}

func (s *activityService) CompleteActivity(ctx context.Context, id uuid.UUID) error {
	activityRealization, err := s.repo.GetRealizationByID(ctx, id)
	if err != nil {
		return err
	}

	if activityRealization.Status != domain.StatusInProgress {
		return fmt.Errorf("Cannot complete activity: current status is %s", activityRealization.Status)
	}

	now := time.Now()
	activityRealization.Status = domain.StatusCompleted
	activityRealization.FinishedAt = &now

	return s.repo.UpdateRealization(ctx, activityRealization)
}
