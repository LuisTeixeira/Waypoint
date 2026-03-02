package handler

import (
	"html/template"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/luisteixeira/waypoint/backend/internal/domain"
	"github.com/luisteixeira/waypoint/backend/internal/ui"
)

type UIHandler struct {
	service domain.ActivityService
	tmpl    *template.Template
}

func NewUIHandler(svc domain.ActivityService) *UIHandler {
	tmpl := template.Must(template.ParseFS(ui.Files, "templates/*.html", "templates/partials/*.html"))
	return &UIHandler{service: svc, tmpl: tmpl}
}

func (h *UIHandler) ShowDashboard(w http.ResponseWriter, r *http.Request) {
	// For testing this will come from auth session
	familyID := uuid.MustParse(os.Getenv("TEST_FAMILY_ID"))
	entityID := uuid.MustParse(os.Getenv("TEST_ENTITY_ID"))
	data := map[string]interface{}{
		"TestFamilyID": familyID,
		"TestEntityID": entityID,
	}
	h.tmpl.ExecuteTemplate(w, "layout.html", data)
}
