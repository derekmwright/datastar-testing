package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/starfederation/datastar-go/datastar"

	"exampleapp/internal/views"
)

func LandingPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Full page reload
		if r.Header.Get("Datastar-Request") != "true" {
			Root(r.URL.Path)(w, r)
			return
		}

		// Otherwise we have a datastar request; upgrade the connection to SSE and Patch elements and update navigation
		sse := datastar.NewSSE(w, r)

		if err := sse.PatchElementGostar(views.LandingPage()); err != nil {
			log.Printf("Unable to patch landing page, error: %v", err)
		}

		if err := sse.ExecuteScript(
			fmt.Sprintf("history.pushState({}, '', '%s');", r.URL.Path),
		); err != nil {
			log.Printf("Unable to send SSE, error: %v", err)
		}
	}
}
