package main

import (
	"github.com/dev-ayaa/resvbooking/pkg/config"
	"github.com/go-chi/chi"
	"testing"
)

func TestRoutes(t *testing.T) {

	var app *config.AppConfig
	mux := routes(app)
	switch mux.(type) {
	case *chi.Mux:
		//test successful
	default:
		t.Errorf("Testing for Routes .....\n%vIs not a Chi httpHandler....", mux)
	}
}
