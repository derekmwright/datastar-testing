package main

import (
	. "github.com/derekmwright/htemel"
	. "github.com/derekmwright/htemel/html"
)

// Site takes a targetView string which informs datastar which view to request from the backend.
// This view should only be called during full page reloads.
func Site(targetView string) Node {
	return Group(
		GenericVoid("!DOCTYPE", map[string]any{"html": nil}),
		Html(
			Head(
				Meta().Charset("utf-8"),
				Meta().Name("viewport").Content("width=device-width, initial-scale=1"),
				Title(Text("Data-Star Testing")),
				Script().Type("module").Src("https://cdn.jsdelivr.net/gh/starfederation/datastar@main/bundles/datastar.js"),
				Script().Src("https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"),
			),
			Body(
				Div().Id("app-view").Data("on-load", "@get('"+targetView+"')"),
			).Class("text-gray-200"),
		).Id("page-root").Lang("en").Class("h-dvh bg-gray-900"),
	)
}

// LandingPage is the default page that a user is shown when navigating to the site.
func LandingPage() Node {
	return Div(
		SiteNav("landing-page"),
		H1(Text("Welcome to the home page")).Class("text-xl font-semibold"),
	).Id("app-view")
}

// UserProfile shows some user profile information
func UserProfile() Node {
	return Div(
		SiteNav("user-profile"),
		H1(Text("User profile")).Class("text-xl font-semibold"),
	).Id("app-view")
}

func SiteNav(activeUrl string) Node {
	return Div(
		Nav(
			Ul(
				NavLink("Home", "/landing-page", activeUrl == "landing-page"),
				NavLink("User Profile", "/user-profile", activeUrl == "user-profile"),
			),
		),
	).Id("navigation-container")
}

func NavLink(name, url string, active bool) Node {
	classes := "hover:text-gray-300 hover:border-b hover:border-b-gray-300"
	if active {
		classes += "text-gray-300 border-b border-b-gray-300"
	}

	return Li(
		A(
			Text(name),
		).
			Href(url).
			Class(classes).
			Data("on-click__prevent", "@get('"+url+"')"),
	)
}
