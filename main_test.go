package main

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

// The iOS app depends on these routes. See AGENTS.md before changing them.
func TestRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	got := map[string]bool{}
	for _, route := range initRouter().Routes() {
		got[route.Method+" "+route.Path] = true
	}

	for _, want := range []string{
		http.MethodGet + " /",
		http.MethodGet + " /sheets/default",
		http.MethodGet + " /oauth/google",
		http.MethodGet + " /oauth/google/callback",
		http.MethodGet + " /sheets",
		http.MethodPost + " /sheets",
		http.MethodGet + " /sheets/:id",
		http.MethodDelete + " /sheets/:id",
	} {
		if !got[want] {
			t.Errorf("route %q isn't registered", want)
		}
	}
}
