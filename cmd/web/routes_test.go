package main

import (
	"fmt"
	"github.com/dev-ayaa/resvbooking/pkg/config"
	"github.com/go-chi/chi"
	"testing"
)

func TestRoutes(t *testing.T) {

	var app config.AppConfig
	mux := routes(&app)
	switch rt := mux.(type) {
	case *chi.Mux:
		//test successful
	default:
		t.Error(fmt.Sprintf("Testing for Routes .....\n%TIs not a Chi httpHandler....", rt))
	}
}
