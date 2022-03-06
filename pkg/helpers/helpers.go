package helpers

import (
	"fmt"
	"github.com/dev-ayaa/resvbooking/pkg/config"
	"net/http"
	"runtime/debug"
)

var app *config.AppConfig

func NewHelper(a *config.AppConfig) {
	app = a
}

func ClientSideError(wr http.ResponseWriter, statusCode int) {
	app.InfoLog.Println("Client Error with the status code ", statusCode)
	http.Error(wr, http.StatusText(statusCode), statusCode)
}

func ServerSideError(wr http.ResponseWriter, err error) {
	trackedError := fmt.Sprintf("%v....\n%v....", err.Error(), debug.Stack())
	app.ErrorLog.Println(trackedError)
	http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
