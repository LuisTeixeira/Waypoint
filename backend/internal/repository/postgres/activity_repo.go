package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/luisteixeira/waypoint/backend/internal/domain"
	"github.com/luisteixeira/waypoint/backend/internal/repository"
)

type postgresActivityRepo struct {
	db *sql.DB
}

func NewPostgresActivityRepo(db *sql.DB) *postgresActivityRepo {
	return &postgresActivityRepo{db: db}
}

func (r *postgresActivityRepo) CreateRealization(ctx context.Context, activityRealization *domain.ActivityRealization) error {
	familyID, err := repository.GetFamilyIdFromContext(ctx)
	if err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, `
		INSERT INTO activity_realizations (family_id, definition_id, entity_id, status, started_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		familyID, activityRealization.DefinitionID, activityRealization.EntityID, activityRealization.Status, activityRealization.StartedAt,
	).Scan(&activityRealization.ID)
	if err != nil {
		return err
	}

	for _, caregiverID := range activityRealization.CaregiversIDs {
		_, err := tx.ExecContext(ctx,
			"INSERT INTO realization_caregivers (realization_id, caregiver_id) VALUES ($1, $2)",
			activityRealization.ID, caregiverID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *postgresActivityRepo) GetRealizationByID(ctx context.Context, id uuid.UUID) (*domain.ActivityRealization, error) {
	familyID, err := repository.GetFamilyIdFromContext(ctx)
	if err != nil {
		return nil, err
	}

	query := `
			SELECT
				ar.id, ar.family_id, ar.definition_id, ar.entity_id, ar.status,
				ar.started_at, ar.finished_at,
				COALESCE(array_agg(rc.caregiver_id) FILTER(WHERE rc.caregiver_id IS NOT NULL), '{}') as caregiver_ids
			FROM activity_realizations ar
			LEFT JOIN realization_caregivers rc ON ar.id = rc.realization_id
			WHERE ar.id = $1 AND ar.family_id = $2
			GROUP BY ar.id;`

	var activity_realization domain.ActivityRealization
	var caregiverIDs []uuid.UUID

	err = r.db.QueryRowContext(ctx, query, id, familyID).Scan(
		&activity_realization.ID,
		&activity_realization.FamilyID,
		&activity_realization.DefinitionID,
		&activity_realization.EntityID,
		&activity_realization.Status,
		&activity_realization.StartedAt,
		&activity_realization.FinishedAt,
		pq.Array(&caregiverIDs),
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Activity realization not found")
		}
		return nil, fmt.Errorf("Failed to fetch realization: %w", err)
	}

	activity_realization.CaregiversIDs = caregiverIDs
	return &activity_realization, nil
}

func (r *postgresActivityRepo) GetActiveByEntity(ctx context.Context, entityID uuid.UUID) (*domain.ActivityRealization, error) {
	familyID, err := repository.GetFamilyIdFromContext(ctx)
	if err != nil {
		return nil, err
	}

	query := `
			SELECT
				ar.id, ar.family_id, ar.definition_id, ar.entity_id, ar.status,
				ar.started_at, ar.finished_at,
				COALESCE(array_agg(rc.caregiver_id) FILTER (WHERE rc.caregiver_id IS NOT NULL), '{}')
			FROM activity_realizations as ar
			LEFT JOIN realization_caregivers rc ON ar.id = rc.realization_id
			WHERE ar.entity_id = $1 AND ar.family_id = $2 AND ar.status = $3
			GROUP bY ar.id
			LIMIT 1;
	`

	var ar domain.ActivityRealization
	var caregiverIDs []uuid.UUID

	err = r.db.QueryRowContext(ctx, query, entityID, familyID, domain.StatusInProgress).Scan(
		&ar.ID, &ar.FamilyID, &ar.DefinitionID, &ar.EntityID, &ar.Status,
		&ar.StartedAt, &ar.FinishedAt, pq.Array(&caregiverIDs),
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to check active status: %w", err)
	}

	ar.CaregiversIDs = caregiverIDs
	return &ar, nil
}

func (r *postgresActivityRepo) UpdateRealization(ctx context.Context, activityRealization *domain.ActivityRealization) error {
	familyID, err := repository.GetFamilyIdFromContext(ctx)
	if err != nil {
		return err
	}

	query := `
			UPDATE activity_realizations
			SET status = $1, started_at=$2, finished_at = $3
			WHERE id = $4 and family_id = $5
	`

	_, err = r.db.ExecContext(ctx, query, activityRealization.Status, activityRealization.StartedAt, activityRealization.FinishedAt,
		activityRealization.ID, familyID)
	return err
}
