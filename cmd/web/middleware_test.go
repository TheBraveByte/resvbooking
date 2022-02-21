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
		t.Errorf("Testing for NoSurf .....\n%v is not http Handler....", ts)
	}

}

func TestSessionLoad(t *testing.T) {
	var h *ResvHandler
	tsl := SessionLoad(h)

	switch tsl.(type) {
	case http.Handler:
		break

	default:
		t.Errorf("Testing for SessionLoad .....\n%v is not http Handler....", tsl)
	}
}
