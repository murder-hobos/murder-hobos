package routes

import (
	"context"
	"net/http"
)

func withSourceIDs(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sourceIDs := []string{"1", "2", "3"}
		ctx := context.WithValue(r.Context(), "sourceIDs", sourceIDs)
		r = r.WithContext(ctx)
		fn(w, r)
	}
}
