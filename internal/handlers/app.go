package handlers

import (
	"fmt"
	"net/http"

	"dstartest/internal/views"
)

func Health(ready bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		if _, err := fmt.Fprint(w, ready); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func Root(page string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := views.Site(page).Render(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
