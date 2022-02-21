package main

import (
	"net/http"
	"os"
	"testing"
)

//setting up our http.Handler interfaces

type ResvHandler struct {
}

func TestMain(m *testing.M) {

	os.Exit(m.Run())
}

func (rs *ResvHandler) ServeHTTP(wr http.ResponseWriter, rq *http.Request) {

}
