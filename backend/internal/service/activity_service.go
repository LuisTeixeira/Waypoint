package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/luisteixeira/waypoint/backend/internal/domain"
)

type activityService struct {
	repo    ActivityRepository
	defRepo DefinitionRepository
}

func NewActivityService(repo ActivityRepository, defRepo DefinitionRepository) *activityService {
	return &activityService{
		repo:    repo,
		defRepo: defRepo,
	}
}

func (s *activityService) StartActivity(ctx context.Context, input domain.StartActivityInput) (*domain.ActivityRealization, error) {
	var realization *domain.ActivityRealization
	var err error

	if input.RealizationID != uuid.Nil {
		realization, err = s.repo.GetRealizationByID(ctx, input.RealizationID)
		if err != nil {
			return nil, err
		}

		if realization.Status != domain.StatusPlanned {
			return nil, fmt.Errorf("cannot start activity: current status is %s", realization.Status)
		}
	} else {
		defID, err := s.resolveDefinitionID(ctx, input)
		if err != nil {
			return nil, err
		}

		realization = &domain.ActivityRealization{
			ID:            uuid.New(),
			DefinitionID:  defID,
			EntityID:      input.EntityID,
			CaregiversIDs: input.CaregiversIDs,
			Status:        domain.StatusPlanned,
		}
	}

	active, err := s.repo.GetActiveByEntity(ctx, input.EntityID)
	if err != nil {
		return nil, fmt.Errorf("failed to check child status: %w", err)
	}
	if active != nil {
		return nil, domain.ErrEntityBusy
	}

	now := time.Now()
	realization.Status = domain.StatusInProgress
	realization.StartedAt = &now

	if input.RealizationID != uuid.Nil {
		err = s.repo.UpdateRealization(ctx, realization)
	} else {
		err = s.repo.CreateRealization(ctx, realization)
	}
	if err != nil {
		return nil, err
	}
	return realization, nil
}

func (s *activityService) PlanActivity(ctx context.Context, input domain.StartActivityInput) (*domain.ActivityRealization, error) {
	defID, err := s.resolveDefinitionID(ctx, input)
	if err != nil {
		return nil, err
	}

	activityRealization := &domain.ActivityRealization{
		DefinitionID:  defID,
		EntityID:      input.EntityID,
		CaregiversIDs: input.CaregiversIDs,
		Status:        domain.StatusPlanned,
	}

	if err := s.repo.CreateRealization(ctx, activityRealization); err != nil {
		return nil, err
	}

	return activityRealization, nil
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

func (s *activityService) resolveDefinitionID(ctx context.Context, input domain.StartActivityInput) (uuid.UUID, error) {
	if input.DefinitionID != uuid.Nil {
		return input.DefinitionID, nil
	}

	if input.NewDefinittionName == "" {
		return uuid.Nil, fmt.Errorf("either definition_id or new_definition_name must be provided")
	}

	def, err := s.defRepo.GetOrCreateByName(ctx, input.NewDefinittionName)
	if err != nil {
		return uuid.Nil, err
	}
	return def.ID, nil
}
