package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type postFormData struct {
	formKey   string
	formValue string
}

var hTests = []struct {
	pageName       string
	pagesUrl       string
	pageMethod     string
	pageFormData   postFormData
	pageStatusCode int
}{
	{pageName: "HomePage", pagesUrl: "/", pageMethod: "GET", pageFormData: postFormData{}, pageStatusCode: http.StatusOK},
	{pageName: "JuniorSuitePage", pagesUrl: "/junior-suite", pageMethod: "GET", pageFormData: postFormData{}, pageStatusCode: http.StatusOK},
	{pageName: "DeluxeSuitePage", pagesUrl: "/deluxe-suite", pageMethod: "GET", pageFormData: postFormData{}, pageStatusCode: http.StatusOK},
	{pageName: "MakeReservationPage", pagesUrl: "/make-reservation", pageMethod: "GET", pageFormData: postFormData{}, pageStatusCode: http.StatusOK},
	{pageName: "PostMakeReservationPage", pagesUrl: "/make-reservation", pageMethod: "POST", pageFormData: postFormData{}, pageStatusCode: http.StatusOK},
	{pageName: "MakeReservationSummary", pagesUrl: "/make-reservation-data", pageMethod: "GET", pageFormData: postFormData{}, pageStatusCode: http.StatusOK},
	{pageName: "CheckAvailabilityPage", pagesUrl: "/check-availability", pageMethod: "GET", pageFormData: postFormData{}, pageStatusCode: http.StatusOK},
	{pageName: "PostCheckAvailabilityPage", pagesUrl: "/check-availability", pageMethod: "POST", pageFormData: postFormData{}, pageStatusCode: http.StatusOK},
	{pageName: "AboutPage", pagesUrl: "/about", pageMethod: "GET", pageFormData: postFormData{}, pageStatusCode: http.StatusOK},
	{pageName: "ContactPage", pagesUrl: "/contact", pageMethod: "GET", pageFormData: postFormData{}, pageStatusCode: http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	//web server that returns a status code, a test Server and a client that
	//can ca on the server
	sth := httptest.NewTLSServer(routes)
	defer sth.Close()

	for _, h := range hTests {
		if h.pageMethod == "GET" {
			response, err := sth.Client().Get(sth.URL + h.pagesUrl)
			if err != nil {
				t.Log(err)
				t.Fatal(fmt.Sprintf("Error Testing the handler of %s......%s", h.pageName, err))
			}
			if response.StatusCode != h.pageStatusCode {
				t.Errorf("Error statusCode for %s : get %v statusCode instead of %v", h.pageName, response.StatusCode, h.pageStatusCode)
			}
		} else {

		}
	}

}
