package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/luisteixeira/waypoint/backend/internal/domain"
	"github.com/luisteixeira/waypoint/backend/internal/handler"
	"github.com/luisteixeira/waypoint/backend/internal/middleware"
	"github.com/luisteixeira/waypoint/backend/internal/service"
	"github.com/luisteixeira/waypoint/backend/test/internal/repository/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActivityHandler_PlanAndStart(t *testing.T) {
	router := setupTestRouter()
	familyID := uuid.New().String()
	entityID := uuid.New()

	var plannedID uuid.UUID

	t.Run("Plan Activity", func(t *testing.T) {
		payload := map[string]interface{}{
			"entity_id":           entityID,
			"new_definition_name": "Afternoon nap",
		}
		body, _ := json.Marshal(payload)

		request := httptest.NewRequest("POST", "/api/v1/activities/plan", bytes.NewBuffer(body))
		request.Header.Set("X-Family-ID", familyID)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, request)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response domain.ActivityRealization
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, domain.StatusPlanned, response.Status)
		plannedID = response.ID
	})

	t.Run("Start the planned Activity", func(t *testing.T) {
		payload := map[string]interface{}{
			"realization_id": plannedID,
			"entity_id":      entityID,
		}
		body, _ := json.Marshal(payload)

		request := httptest.NewRequest("POST", "/api/v1/activities/start", bytes.NewBuffer(body))
		request.Header.Set("X-Family-ID", familyID)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, request)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response domain.ActivityRealization
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, domain.StatusInProgress, response.Status)
		assert.NotNil(t, response.StartedAt)
	})

	t.Run("Fail to start another activity when one is already in progress (Conflict)", func(t *testing.T) {
		payload := map[string]interface{}{
			"entity_id":           entityID,
			"new_definition_name": "Sport",
		}
		body, _ := json.Marshal(payload)

		request := httptest.NewRequest("POST", "/api/v1/activities/start", bytes.NewBuffer(body))
		request.Header.Set("X-Family-ID", familyID)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, request)

		assert.Equal(t, http.StatusConflict, w.Code)
	})
}

func setupTestRouter() *chi.Mux {
	activityRepo := memory.NewInMemoryActivityRepo()
	definitionRepo := memory.NewInMemoryDefinitionRepo()
	svc := service.NewActivityService(activityRepo, definitionRepo)
	handler := handler.NewActivityHandler(svc)

	router := chi.NewRouter()
	router.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.TenantMiddleware)
		r.Post("/activities/plan", handler.PlanActivity)
		r.Post("/activities/start", handler.StartActivity)
	})

	return router
}
