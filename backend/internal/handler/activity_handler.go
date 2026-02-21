package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/luisteixeira/waypoint/backend/internal/domain"
)

type CreateActivityRequest struct {
	EntityID           uuid.UUID   `json:"entity_id"`
	DefinitionID       uuid.UUID   `json:"definition_id,omitempty"`
	NewDefinittionName string      `json:"new_definition_name,omitempty"`
	CaregiverIDs       []uuid.UUID `json:"caregiver_ids"`
}

type ActivityHandler struct {
	service domain.ActivityService
}

func NewActivityHandler(service domain.ActivityService) *ActivityHandler {
	return &ActivityHandler{service: service}
}

func (h *ActivityHandler) PlanActivity(w http.ResponseWriter, r *http.Request) {
	var createActivityRequest CreateActivityRequest

	if err := json.NewDecoder(r.Body).Decode(&createActivityRequest); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
	}

	input := domain.StartActivityInput{
		EntityID:           createActivityRequest.EntityID,
		DefinitionID:       createActivityRequest.DefinitionID,
		NewDefinittionName: createActivityRequest.NewDefinittionName,
		CaregiversIDs:      createActivityRequest.CaregiverIDs,
	}

	activityRealization, err := h.service.PlanActivity(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(activityRealization)
}
