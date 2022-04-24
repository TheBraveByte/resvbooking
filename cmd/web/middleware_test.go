package main

import (
	"net/http"
	"testing"
)

//Testing for MiddleWares
func TestNoSurf(t *testing.T) {
	var h *ResvHandler
	ts := NoSurf(h)

	switch ts.(type) {
	case http.Handler:
		break

	default:
		t.Errorf("Testing for NoSurf .....\n%vIs not an http Handler....", ts)
	}
}

func TestSessionLoad(t *testing.T) {
	var h *ResvHandler
	tsl := SessionLoad(h)

	switch tsl.(type) {
	case http.Handler:
		break

	default:
		t.Errorf("Testing for SessionLoad .....\n%v Is not an http Handler....", tsl)
	}
}

func TestAuthenticate(t *testing.T) {
	var auth *ResvHandler
	a := Authenticate(auth)
	switch a.(type) {
	case http.Handler:
		break
	default:
		t.Errorf("Testing for Authentication Function .....\n%v is not a http Handler", a)
	}
}
