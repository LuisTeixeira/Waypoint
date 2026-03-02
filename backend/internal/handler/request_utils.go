package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func decodeRequest(r *http.Request, dst interface{}) error {
	contentType := r.Header.Get("Content-Type")

	if strings.Contains(contentType, "application/json") {
		return json.NewDecoder(r.Body).Decode(dst)
	}

	if err := r.ParseForm(); err != nil {
		return err
	}

	if request, ok := dst.(*ActivityRequest); ok {
		if val := r.FormValue("entity_id"); val != "" {
			id, _ := uuid.Parse(val)
			request.EntityID = id
		}
		if val := r.FormValue("realization_id"); val != "" {
			id, _ := uuid.Parse(val)
			request.RealizationID = &id
		}
		request.NewDefinittionName = r.FormValue("new_definition_name")
	}

	return nil
}
