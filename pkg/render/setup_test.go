package render

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/dev-ayaa/resvbooking/pkg/config"
	"github.com/dev-ayaa/resvbooking/pkg/models"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {
	gob.Register(models.ReservationData{})
	testApp.InProduction = false

	infoLogger := log.New(os.Stdout, "INFO ::\t", log.LstdFlags)
	testApp.InfoLog = infoLogger

	errorLogger := log.New(os.Stdout, "ERROR ::\t", log.LstdFlags|log.Lshortfile)
	testApp.ErrorLog = errorLogger

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = testApp.InProduction

	testApp.Session = session
	app = &testApp

	os.Exit(m.Run())
}

type myResponse struct{}

func (rw *myResponse) Header() http.Header {
	h := http.Header{}
	return h
}

func (rw *myResponse) Write(b []byte) (int, error) {
	byteLength := len(b)
	return byteLength, nil
}

func (wr *myResponse) WriteHeader(statusCode int) {}

/*
type bt uint8
type newStatusCode int
//newStatusCode = http.StatusOk

type myResponse interface {
	NewHeader()
	MyWriter(bt)
	MyWriterHeader(newStatusCode)
}

func NewHeader() map[string][]string {
	h := map[string][]string{}
	return h
}

func MyWriter(b []byte) (int, error) {
	for _, x :=range b{
		return int(x), nil
	}
	return 0, nil
}

func MyWriterHeader(statusCode int){}
*/
