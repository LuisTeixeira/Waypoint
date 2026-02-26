package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
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
		NewDefinittionName: activityRequest.NewDefinittionName,
		CaregiversIDs:      activityRequest.CaregiverIDs,
	}

	if activityRequest.RealizationID != nil {
		input.RealizationID = *activityRequest.RealizationID
	}
	if activityRequest.DefinitionID != nil {
		input.DefinitionID = *activityRequest.DefinitionID
	}

	activityRealization, err := h.service.PlanActivity(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderJSON(w, http.StatusCreated, activityRealization)
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
			renderError(w, err.Error(), http.StatusConflict)
			return
		}

		if err.Error() == "not found" {
			renderError(w, "planned activity not found", http.StatusNotFound)
		}

		log.Printf("StartActivity Error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	renderJSON(w, http.StatusCreated, activityRealization)
}

func (h *ActivityHandler) CompleteActivity(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		renderError(w, "invalid activity id", http.StatusBadRequest)
		return
	}

	err = h.service.CompleteActivity(r.Context(), id)
	if err != nil {
		if err.Error() == "not found" {
			renderError(w, "activity not found", http.StatusNotFound)
			return
		}

		log.Printf("CompletedActivity Error: %v", err)
		renderError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func renderJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func renderError(w http.ResponseWriter, message string, status int) {
	renderJSON(w, status, map[string]string{"error": message})
}
