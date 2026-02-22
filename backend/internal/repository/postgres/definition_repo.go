package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/luisteixeira/waypoint/backend/internal/domain"
	"github.com/luisteixeira/waypoint/backend/internal/repository"
)

type postgresDefinitionRepo struct {
	db *sql.DB
}

func NewPostgresDefinitionRepo(db *sql.DB) *postgresDefinitionRepo {
	return &postgresDefinitionRepo{db: db}
}

func (r *postgresDefinitionRepo) GetOrCreateByName(ctx context.Context, name string) (*domain.ActivityDefinition, error) {
	familyID, err := repository.GetFamilyIdFromContext(ctx)
	if err != nil {
		return nil, err
	}

	query := `
			INSERT INTO activity_definitions (family_id, name)
			VALUES ($1, $2)
			ON CONFLICT (family_id, name) DO UPDATE SET name = EXCLUDED.name
			RETURNING id, family_id, name, description, color_code;
	`

	var def domain.ActivityDefinition
	err = r.db.QueryRowContext(ctx, query, familyID, name).Scan(
		&def.ID, &def.FamilyID, &def.Name, &def.Description, &def.ColorCode,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get or create definition: %w", err)
	}
	return &def, nil
}

func (r *postgresDefinitionRepo) ListByFamily(ctx context.Context) ([]domain.ActivityDefinition, error) {
	familyID, err := repository.GetFamilyIdFromContext(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, family_id, name, description, color_code
		FROM activity_definitions
		WHERE family_id = $1 ORDER BY name ASC`,
		familyID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var defs []domain.ActivityDefinition
	for rows.Next() {
		var d domain.ActivityDefinition
		if err := rows.Scan(&d.ID, &d.FamilyID, &d.Name, &d.Description, &d.ColorCode); err != nil {
			return nil, err
		}
		defs = append(defs, d)
	}
	return defs, nil
}
