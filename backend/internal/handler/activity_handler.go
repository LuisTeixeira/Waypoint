package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/luisteixeira/waypoint/backend/internal/domain"
)

type ActivityRequest struct {
	RealizationID      *uuid.UUID  `json:"realization_id,omitempty"`
	EntityID           uuid.UUID   `json:"entity_id"`
	DefinitionID       *uuid.UUID  `json:"definition_id,omitempty"`
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
	var activityRequest ActivityRequest

	if err := json.NewDecoder(r.Body).Decode(&activityRequest); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
	}

	input := domain.StartActivityInput{
		EntityID:           activityRequest.EntityID,
		DefinitionID:       *activityRequest.DefinitionID,
		NewDefinittionName: activityRequest.NewDefinittionName,
		CaregiversIDs:      activityRequest.CaregiverIDs,
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

func (h *ActivityHandler) StartActivity(w http.ResponseWriter, r *http.Request) {
	var activityRequest ActivityRequest
	if err := json.NewDecoder(r.Body).Decode(&activityRequest); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	input := domain.StartActivityInput{
		EntityID:           activityRequest.EntityID,
		NewDefinittionName: activityRequest.NewDefinittionName,
		CaregiversIDs:      activityRequest.CaregiverIDs,
	}

	if activityRequest.RealizationID != nil {
		input.RealizationID = *activityRequest.RealizationID
	}
	if activityRequest.DefinitionID != nil {
		input.DefinitionID = *activityRequest.DefinitionID
	}

	activityRealization, err := h.service.StartActivity(r.Context(), input)
	if err != nil {
		if errors.Is(err, domain.ErrEntityBusy) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if err.Error() == "not found" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "planned activity not found"})
		}

		log.Printf("StartActivity Error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(activityRealization)
}
