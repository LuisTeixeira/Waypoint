package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const FamilyIDKey contextKey = "family_id"

func TenantMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		familyIDStr := r.Header.Get("X-Family-ID")
		if familyIDStr == "" {
			http.Error(w, "Missing Family Context", http.StatusUnauthorized)
			return
		}

		familyID, err := uuid.Parse(familyIDStr)
		if err != nil {
			http.Error(w, "Invalid Family ID", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), FamilyIDKey, familyID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
