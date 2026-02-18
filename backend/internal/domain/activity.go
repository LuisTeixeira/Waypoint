package domain

import (
	"time"

	"github.com/google/uuid"
)

type ActivityStatus string

const (
	StatusPlanned    ActivityStatus = "planned"
	StatusInProgress ActivityStatus = "in_progress"
	StatusCompleted  ActivityStatus = "completed"
	StatusCancelled  ActivityStatus = "cancelled"
)

type ActivityDefinition struct {
	ID       uuid.UUID `json:"id"`
	FamilyID uuid.UUID `json:"family_id"`
	Name     string    `json:"name"`
}

type ActivityRealization struct {
	ID            uuid.UUID      `json:"id"`
	FamilyID      uuid.UUID      `json:"family_id"`
	DefinitionID  uuid.UUID      `json:"definition_id"`
	EntityID      uuid.UUID      `json:"entity_id"`
	CaregiversIDs []uuid.UUID    `json:"caregiver_ids"`
	Status        ActivityStatus `json:"status"`
	StartedAt     *time.Time     `json:"started_at"`
	FinishedAt    *time.Time     `json:"finished_at"`
}
