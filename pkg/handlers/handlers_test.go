package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postFormData struct {
	formKey   string
	formValue string
}

//Slices of struct
var hTests = []struct {
	pageName       string
	pagesUrl       string
	pageMethod     string
	pageFormData   []postFormData
	pageStatusCode int
}{
	{pageName: "HomePage", pagesUrl: "/", pageMethod: "GET", pageFormData: []postFormData{}, pageStatusCode: http.StatusOK},

	{pageName: "JuniorSuitePage", pagesUrl: "/junior-suite", pageMethod: "GET", pageFormData: []postFormData{}, pageStatusCode: http.StatusOK},

	{pageName: "DeluxeSuitePage", pagesUrl: "/deluxe-suite", pageMethod: "GET", pageFormData: []postFormData{}, pageStatusCode: http.StatusOK},

	{pageName: "MakeReservationPage", pagesUrl: "/make-reservation", pageMethod: "GET", pageFormData: []postFormData{}, pageStatusCode: http.StatusOK},

	{pageName: "PostMakeReservationPage", pagesUrl: "/make-reservation", pageMethod: "POST", pageFormData: []postFormData{
		{formKey: "first-name", formValue: "Yusuf"},
		{formKey: "last-name", formValue: "Abidemi Akinleye"},
		{formKey: "email", formValue: "YusufAkinleyeAbidemi@gmail.com"},
		{formKey: "phone-number", formValue: "+23409096000258"},
		{formKey: "inputPassword", formValue: "***************"},
		{formKey: "inputPassword4", formValue: "***************"},
	}, pageStatusCode: http.StatusOK},

	{pageName: "MakeReservationSummary", pagesUrl: "/make-reservation-data", pageMethod: "GET", pageFormData: []postFormData{}, pageStatusCode: http.StatusOK},

	{pageName: "CheckAvailabilityPage", pagesUrl: "/check-availability", pageMethod: "GET", pageFormData: []postFormData{}, pageStatusCode: http.StatusOK},

	{pageName: "PostCheckAvailabilityPage", pagesUrl: "/check-availability", pageMethod: "POST", pageFormData: []postFormData{
		{formKey: "check-in", formValue: "2000-27-01"},
		{formKey: "check-in", formValue: "2000-24-02"},
	},
		pageStatusCode: http.StatusOK},

	{pageName: "AboutPage", pagesUrl: "/about", pageMethod: "GET", pageFormData: []postFormData{}, pageStatusCode: http.StatusOK},

	{pageName: "ContactPage", pagesUrl: "/contact", pageMethod: "GET", pageFormData: []postFormData{}, pageStatusCode: http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	//web server that returns a status code, a test Server and a client must be setup
	sth := httptest.NewTLSServer(routes)

	//close the server after the main function execute
	defer sth.Close()

	for _, h := range hTests {
		if h.pageMethod == "GET" {
			//this is when the client is requesting to view a page
			//server.URL == localhost
			//h.pagesUrls is the page the client is requesting for
			response, err := sth.Client().Get(sth.URL + h.pagesUrl)
			if err != nil {
				t.Log(err)
				t.Fatal(fmt.Sprintf("Error Testing the handler of %s......%s", h.pageName, err))
			}
			if response.StatusCode != h.pageStatusCode {
				t.Errorf("Error statusCode for %s : get %v statusCode instead of %v", h.pageName, response.StatusCode, h.pageStatusCode)
			}
		} else {
			//i.e if it is a POST request
			formData := url.Values{}
			for _, v := range h.pageFormData {
				//formData[v.formKey] = []string{v.formValue}
				formData.Add(v.formKey, v.formValue)
			}
			response, err := sth.Client().PostForm(sth.URL+h.pagesUrl, formData)
			if err != nil {
				t.Log(err)
				t.Fatal(fmt.Sprintf("Error Testing the handler of %s......%s", h.pageName, err))
			}
			if response.StatusCode != h.pageStatusCode {
				t.Errorf("Error statusCode for %s : get %v statusCode instead of %v", h.pageName, response.StatusCode, h.pageStatusCode)
			}
		}
	}

}
