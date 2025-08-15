package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/starfederation/datastar-go/datastar"
)

const (
	LandingPageName     string = "landing-page"
	UserProfilePageName string = "user-profile"
)

func rootHandler(page string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		if err := Site(page).Render(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func landingHandler(w http.ResponseWriter, r *http.Request) {
	// Full page reload
	if r.Header.Get("Datastar-Request") != "true" {
		rootHandler(LandingPageName)(w, r)
		return
	}

	// Otherwise we have a datastar request; upgrade the connection to SSE and Patch elements and update navigation
	sse := datastar.NewSSE(w, r)

	if err := sse.PatchElementGostar(LandingPage()); err != nil {
		log.Printf("Unable to patch landing page, error: %v", err)
	}

	if err := sse.ExecuteScript(`history.pushState({}, "Test", "/` + LandingPageName + `");`); err != nil {
		log.Printf("Unable to send SSE, error: %v", err)
	}
}

func userProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Full page reload
	if r.Header.Get("Datastar-Request") != "true" {
		rootHandler(UserProfilePageName)(w, r)
		return
	}

	// Otherwise we have a datastar request; upgrade the connection to SSE and Patch elements and update navigation
	sse := datastar.NewSSE(w, r)

	if err := sse.PatchElementGostar(UserProfile()); err != nil {
		log.Printf("Unable to patch user profile, error: %v", err)
	}

	if err := sse.ExecuteScript(`history.pushState({}, "Test", "/` + UserProfilePageName + `");`); err != nil {
		log.Printf("Unable to send SSE, error: %v", err)
	}
}

func main() {
	r := chi.NewRouter()
	r.Get("/", rootHandler(LandingPageName))
	r.Get("/landing-page", landingHandler)
	r.Get("/user-profile", userProfileHandler)
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
