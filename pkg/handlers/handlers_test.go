package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/dev-ayaa/resvbooking/pkg/driver"
	"github.com/dev-ayaa/resvbooking/pkg/models"
)

type postFormData struct {
	formKey   string
	formValue string
}

//var session *scs.SessionManager

//Slices of struct
var hTests = []struct {
	pageName       string
	pagesUrl       string
	pageMethod     string
	pageStatusCode int
}{
	{pageName: "HomePage", pagesUrl: "/", pageMethod: "GET", pageStatusCode: http.StatusOK},

	{pageName: "JuniorSuitePage", pagesUrl: "/junior-suite", pageMethod: "GET", pageStatusCode: http.StatusOK},

	{pageName: "DeluxeSuitePage", pagesUrl: "/deluxe-suite", pageMethod: "GET", pageStatusCode: http.StatusOK},

	// {pageName: "MakeReservationPage", pagesUrl: "/make-reservation", pageMethod: "GET", pageStatusCode: http.StatusOK},

	// {pageName: "MakeReservationSummary", pagesUrl: "/make-reservation-data", pageMethod: "GET", pageStatusCode: http.StatusOK},

	// {pageName: "CheckAvailabilityPage", pagesUrl: "/check-availability", pageMethod: "GET", pageStatusCode: http.StatusOK},
	{pageName: "AboutPage", pagesUrl: "/about", pageMethod: "GET", pageStatusCode: http.StatusOK},

	{pageName: "ContactPage", pagesUrl: "/contact", pageMethod: "GET", pageStatusCode: http.StatusOK},
}

//Notice
// create our request with a nil body, so parsing form fails
// rq, _ = http.NewRequest("POST", "/search-availability", nil)

// // get the context with session
// ctx = getContext(rq)
// rq = rq.WithContext(ctx)

// // set the request header
// rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

// // create our response recorder, which satisfies the requirements
// // for http.ResponseWriter
// responseRecorder = httptest.NewRecorder()

// // make our handler a http.HandlerFunc
// handler = http.HandlerFunc(Repo.PostCheckAvailabilityPage)

// // make the request to our handler
// handler.ServeHTTP(responseRecorder, rq)

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	//web server that returns a status code, a test Server and a client must be setup
	sth := httptest.NewTLSServer(routes)

	//close the server after the main function execute
	defer sth.Close()

	for _, h := range hTests {
		//if h.pageMethod == "GET" {
		//	//this is when the client is requesting to view a page
		//	//server.URL == localhost
		//	//h.pagesUrls is the page the client is requesting for
		response, err := sth.Client().Get(sth.URL + h.pagesUrl)
		if err != nil {
			t.Log(err)
			t.Fatal(fmt.Sprintf("Error Testing the handler of %s......%s", h.pageName, err))
		}
		if response.StatusCode != h.pageStatusCode {
			t.Errorf("Error statusCode for %s : get %v statusCode instead of %v", h.pageName, response.StatusCode, h.pageStatusCode)
		}
		//} else {
		//	//i.e if it is a POST request
		//	formData := url.Values{}
		//	for _, v := range h.pageFormData {
		//		//formData[v.formKey] = []string{v.formValue}
		//		formData.Add(v.formKey, v.formValue)
		//	}
		//	response, err := sth.Client().PostForm(sth.URL+h.pagesUrl, formData)
		//	if err != nil {
		//		t.Log(err)
		//		t.Fatal(fmt.Sprintf("Error Testing the handler of %s......%s", h.pageName, err))
		//	}
		//	if response.StatusCode != h.pageStatusCode {
		//		t.Errorf("Error statusCode for %s : get %v statusCode instead of %v", h.pageName, response.StatusCode, h.pageStatusCode)
		//	}
		//}
	}

}

func TestRepository_MakeReservationPage(t *testing.T) {
	resvTest := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "Deluxe suite",
		},
	}
	rq, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getContext(rq)
	rq = rq.WithContext(ctx)

	responseRecorder := httptest.NewRecorder()
	session.Put(ctx, "reservation", resvTest)
	handler := http.HandlerFunc(Repo.MakeReservationPage)
	handler.ServeHTTP(responseRecorder, rq)

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("Wrong response from the make-reservation handler: got %v wanted %v", responseRecorder.Code, http.StatusOK)
	}

	//Testing when Reservation is not put in session
	rq, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getContext(rq)
	rq = rq.WithContext(ctx)

	responseRecorder = httptest.NewRecorder()
	session.Put(ctx, "reservation", resvTest)
	handler.ServeHTTP(responseRecorder, rq)

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("Wrong response from the make-reservation handler: got %v wanted %v", responseRecorder.Code, http.StatusOK)
	}

	//Testing when the room_id of the reserved room does'nt exist
	rq, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getContext(rq)
	rq = rq.WithContext(ctx)
	resvTest.RoomID = 8
	responseRecorder = httptest.NewRecorder()

	session.Put(ctx, "reservation", resvTest)
	handler.ServeHTTP(responseRecorder, rq)

	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("Wrong response from the make-reservation handler: got %v wanted %v", responseRecorder.Code, http.StatusTemporaryRedirect)
	}

}

func TestRepository_PostMakeReservationPage(t *testing.T) {

	PostRqData := url.Values{}
	PostRqData.Add("check-in", "2022-03-01")
	PostRqData.Add("check-out", "2022-03-06")
	PostRqData.Add("first-name", "Graham")
	PostRqData.Add("last-name", "Graham")
	PostRqData.Add("email", "Grahams@gmail.com")
	PostRqData.Add("phone-number", "20229028844")
	PostRqData.Add("room_id", "1")

	rq, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(PostRqData.Encode()))
	ctx := getContext(rq)
	rq = rq.WithContext(ctx)
	rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostMakeReservationPage)
	handler.ServeHTTP(responseRecorder, rq)

	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("Wrong response from the make-reservation handler: got %v wanted %v", responseRecorder.Code, http.StatusTemporaryRedirect)
	}

	//TESTING FOR INVALID INPUT DATA

	//Tesing for invalid check-in date
	PostRqData.Add("check-in", "invalid")
	PostRqData.Add("check-out", "2022-09-10")
	PostRqData.Add("first-name", "Graham")
	PostRqData.Add("last-name", "Graham")
	PostRqData.Add("email", "Grahams@gmail.com")
	PostRqData.Add("phone-number", "20229028844")
	PostRqData.Add("room_id", "1")

	rq, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(PostRqData.Encode()))
	ctx = getContext(rq)
	rq = rq.WithContext(ctx)
	rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	responseRecorder = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostMakeReservationPage)
	handler.ServeHTTP(responseRecorder, rq)

	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post-resrvation response got invalid data for check-in date: got %v wanted %v", responseRecorder.Code, http.StatusTemporaryRedirect)
	}

	//Tesing for invalid  check-out date
	PostRqData.Add("check-in", "2022-09-09")
	PostRqData.Add("check-out", "invalid")
	PostRqData.Add("first-name", "Graham")
	PostRqData.Add("last-name", "Graham")
	PostRqData.Add("email", "Grahams@gmail.com")
	PostRqData.Add("phone-number", "20229028844")
	PostRqData.Add("room_id", "1")
	rq, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(PostRqData.Encode()))
	ctx = getContext(rq)
	rq = rq.WithContext(ctx)
	rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	responseRecorder = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostMakeReservationPage)
	handler.ServeHTTP(responseRecorder, rq)

	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post-resrvation response got invalid data for check-out date: got %v wanted %v", responseRecorder.Code, http.StatusTemporaryRedirect)
	}

	//Testing for invalid room_id
	PostRqData.Add("check-in", "2022-09-09")
	PostRqData.Add("check-out", "2022-09-10")
	PostRqData.Add("first-name", "Graham")
	PostRqData.Add("last-name", "Graham")
	PostRqData.Add("email", "Grahams@gmail.com")
	PostRqData.Add("phone-number", "20229028844")
	PostRqData.Add("room_id", "invalid")

	rq, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(PostRqData.Encode()))
	ctx = getContext(rq)
	rq = rq.WithContext(ctx)
	rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	responseRecorder = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostMakeReservationPage)
	handler.ServeHTTP(responseRecorder, rq)

	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post-resrvation response got invalid data for room-id date: got %v wanted %v", responseRecorder.Code, http.StatusTemporaryRedirect)

	}

	//Testing for invalid input data
	PostRqData.Add("check-in", "2022-09-09")
	PostRqData.Add("check-out", "2022-09-10")
	PostRqData.Add("first-name", "G")
	PostRqData.Add("last-name", "G")
	PostRqData.Add("email", "Grahams@gmail.com")
	PostRqData.Add("phone-number", "20229028844")
	PostRqData.Add("room_id", "1")

	rq, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(PostRqData.Encode()))
	ctx = getContext(rq)
	rq = rq.WithContext(ctx)
	rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	responseRecorder = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostMakeReservationPage)
	handler.ServeHTTP(responseRecorder, rq)

	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post-resrvation response got invalid data for reservation: got %v wanted %v", responseRecorder.Code, http.StatusTemporaryRedirect)
	}

	// //Tesing for invalid Input data
	// rqField = "check-in=2022-09-09"
	// rqField = fmt.Sprintf("%s&%s", rqField, "check-out=2022-09-10")
	// rqField = fmt.Sprintf("%s&%s", rqField, "first-name=Graham")
	// rqField = fmt.Sprintf("%s&%s", rqField, "last-name=G")
	// rqField = fmt.Sprintf("%s&%s", rqField, "email=Graham@gmail.com")
	// rqField = fmt.Sprintf("%s&%s", rqField, "phone-number=+2349087324561")
	// rqField = fmt.Sprintf("%s&%s", rqField, "room_id=1")

	// rq, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(rqField))
	// ctx = getContext(rq)
	// rq = rq.WithContext(ctx)
	// rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// responseRecorder = httptest.NewRecorder()
	// handler = http.HandlerFunc(Repo.PostMakeReservationPage)
	// handler.ServeHTTP(responseRecorder, rq)

	// if responseRecorder.Code != http.StatusTemporaryRedirect {
	// 	t.Errorf("Wrong response from the make-reservation handler: got %v wanted %v", responseRecorder.Code, http.StatusTemporaryRedirect)
	// }

	//Tesing for insertResevration for database
	PostRqData.Add("check-in", "2022-09-09")
	PostRqData.Add("check-out", "2022-09-10")
	PostRqData.Add("first-name", "Graham")
	PostRqData.Add("last-name", "Graham")
	PostRqData.Add("email", "Grahams@gmail.com")
	PostRqData.Add("phone-number", "20229028844")
	PostRqData.Add("room_id", "14")

	rq, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(PostRqData.Encode()))
	ctx = getContext(rq)
	rq = rq.WithContext(ctx)
	rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	responseRecorder = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostMakeReservationPage)
	handler.ServeHTTP(responseRecorder, rq)

	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post-make reservation failed to insert reservation in the database: got %v wanted %v", responseRecorder.Code, http.StatusTemporaryRedirect)
	}

	//Tesing for InsertRoomRestriction function
	PostRqData.Add("check-in", "2022-09-09")
	PostRqData.Add("check-out", "2022-09-10")
	PostRqData.Add("first-name", "Graham")
	PostRqData.Add("last-name", "Graham")
	PostRqData.Add("email", "Grahams@gmail.com")
	PostRqData.Add("phone-number", "20229028844")
	PostRqData.Add("room_id", "11")
	rq, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(PostRqData.Encode()))
	ctx = getContext(rq)
	rq = rq.WithContext(ctx)
	rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	responseRecorder = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostMakeReservationPage)
	handler.ServeHTTP(responseRecorder, rq)

	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post-Make-Reservation failed to insert reservation because it does not exist in the database : got %v wanted %v", responseRecorder.Code, http.StatusTemporaryRedirect)
	}

}
func TestNewRepo(t *testing.T) {
	var db driver.DB
	testNewRepo := NewRepository(&app, &db)
	if reflect.TypeOf(testNewRepo).String() != "*handlers.Repository" {
		t.Errorf("Did not get the correct repository : got %v wanted *Repository", reflect.TypeOf(NewTestRepository).String())
	}
}

func TestRepository_PostCheckAvailabilityPage(t *testing.T) {
	PostRqData := url.Values{}

	PostRqData.Add("check-in", "2022-09-09")
	PostRqData.Add("check-out", "2022-09-10")

	rq, _ := http.NewRequest("POST", "/check-availability", strings.NewReader(PostRqData.Encode()))
	ctx := getContext(rq)
	rq = rq.WithContext(ctx)
	rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostCheckAvailabilityPage)
	handler.ServeHTTP(responseRecorder, rq)
	if responseRecorder.Code != http.StatusOK {
		t.Errorf("Post availability when no rooms available gave wrong status code: got %d, wanted %d", responseRecorder.Code, http.StatusOK)
	}
	//  rooms are available that does not exist
	PostRqData.Add("check-in", "2000-09-09")
	PostRqData.Add("check-out", "2000-09-10")

	rq, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(PostRqData.Encode()))

	ctx = getContext(rq)
	rq = rq.WithContext(ctx)

	rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostCheckAvailabilityPage)
	handler.ServeHTTP(responseRecorder, rq)

	// since we have rooms available, we expect to get status http.StatusOK
	if responseRecorder.Code != http.StatusOK {
		t.Errorf("Post availability when rooms are available gave wrong status code: got %d, wanted %d", responseRecorder.Code, http.StatusOK)
	}

	//No date to check for reservation
	rq, _ = http.NewRequest("POST", "/search-availability", nil)

	ctx = getContext(rq)
	rq = rq.WithContext(ctx)

	rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	responseRecorder = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostCheckAvailabilityPage)

	handler.ServeHTTP(responseRecorder, rq)

	// since we have rooms available, we expect to get status http.StatusTemporaryRedirect
	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post availability with empty request body (nil) gave wrong status code: got %d, wanted %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	}

	//Wrong  Date format
	// start date in the wrong format
	PostRqData.Add("check-in", "invalid")
	PostRqData.Add("check-out", "2022-09-09")

	rq, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(PostRqData.Encode()))

	// get the context with session
	ctx = getContext(rq)
	rq = rq.WithContext(ctx)

	// set the request header
	rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	responseRecorder = httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler = http.HandlerFunc(Repo.PostCheckAvailabilityPage)

	// make the request to our handler
	handler.ServeHTTP(responseRecorder, rq)

	// since we have rooms available, we expect to get status http.StatusTemporaryRedirect
	if responseRecorder.Code != http.StatusOK {
		t.Errorf("Post availability with invalid start date gave wrong status code: got %d, wanted %d", responseRecorder.Code, http.StatusOK)
	}

	//Wrong date format for check-out date
	PostRqData.Add("check-in", "2022-09-09")
	PostRqData.Add("check-out", "invalid")
	rq, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(PostRqData.Encode()))

	// get the context with session
	ctx = getContext(rq)
	rq = rq.WithContext(ctx)

	// set the request header
	rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	responseRecorder = httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler = http.HandlerFunc(Repo.PostCheckAvailabilityPage)

	// make the request to our handler
	handler.ServeHTTP(responseRecorder, rq)

	// since we have rooms available, we expect to get status http.StatusTemporaryRedirect
	if responseRecorder.Code != http.StatusOK {
		t.Errorf("Post availability with invalid end date gave wrong status code: got %d, wanted %d", responseRecorder.Code, http.StatusOK)
	}

	// database Error fails
	PostRqData.Add("check-in", "2029-09-09")
	PostRqData.Add("check-out", "2029-09-10")
	rq, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(PostRqData.Encode()))

	ctx = getContext(rq)
	rq = rq.WithContext(ctx)
	rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostCheckAvailabilityPage)

	handler.ServeHTTP(responseRecorder, rq)

	// since we have rooms available, we expect to get status http.StatusTemporaryRedirect
	if responseRecorder.Code != http.StatusOK {
		t.Errorf("Post availability when database query fails gave wrong status code: got %d, wanted %d", responseRecorder.Code, http.StatusOK)
	}
}

func TestRepository_AvailabilityJSON(t *testing.T) {
	//No available rooms
	PostRqData := url.Values{}
	PostRqData.Add("check-in", "2050-09-09")
	PostRqData.Add("check-out", "2050-09-10")
	PostRqData.Add("room_id", "1")

	rq, _ := http.NewRequest("POST", "/json-availability", strings.NewReader(PostRqData.Encode()))
	ctx := getContext(rq)
	rq = rq.WithContext(ctx)

	rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.JsonAvailabilityPage)

	handler.ServeHTTP(responseRecorder, rq)

	// since we have no rooms available, we expect to get status http.StatusSeeOther
	// this time we want to parse JSON and get the expected response
	var js ResponseJSON
	err := json.Unmarshal([]byte(responseRecorder.Body.String()), &js)
	if err != nil {
		t.Error("failed to parse json!")
	}

	// since we specified a start date > 2029-09-09, we expect no availability
	if js.Ok {
		t.Error("Got availability when none was expected in AvailabilityJSON")
	}

	/*****************************************
	// second case -- rooms not available
	*****************************************/
	// create our request body
	PostRqData.Add("check-in", "2040-03-01")
	PostRqData.Add("check-out", "2040-03-05")
	PostRqData.Add("room_id", "1")

	rq, _ = http.NewRequest("POST", "/json-availability", strings.NewReader(PostRqData.Encode()))

	ctx = getContext(rq)
	rq = rq.WithContext(ctx)
	rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.JsonAvailabilityPage)

	handler.ServeHTTP(responseRecorder, rq)

	err = json.Unmarshal([]byte(responseRecorder.Body.String()), &js)
	if err != nil {
		t.Error("failed to parse json!")
	}

	// since we specified a start date < 2022-09-09, we expect availability
	if !js.Ok {
		t.Error("Got no availability when some was expected in AvailabilityJSON")
	}

	//No request for PostRqData
	rq, _ = http.NewRequest("POST", "/json-availability", nil)

	ctx = getContext(rq)
	rq = rq.WithContext(ctx)
	rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.JsonAvailabilityPage)

	handler.ServeHTTP(responseRecorder, rq)

	// this time we want to parse JSON and get the expected response
	err = json.Unmarshal([]byte(responseRecorder.Body.String()), &js)
	if err != nil {
		t.Error("failed to parse json!")
	}

	// since we specified a start date < 2022-08-10, we expect availability
	if js.Ok || js.Message != "server error" {
		t.Error("Got availability when request body was empty")
	}

	// make our handler a http.HandlerFunc
	PostRqData.Add("check-in", "2045-09-09")
	PostRqData.Add("check-out", "2045-09-10")
	PostRqData.Add("room_id", "1")

	rq, _ = http.NewRequest("POST", "/json-availability", strings.NewReader(PostRqData.Encode()))

	ctx = getContext(rq)
	rq = rq.WithContext(ctx)

	rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.JsonAvailabilityPage)

	handler.ServeHTTP(responseRecorder, rq)

	// this time we want to parse JSON and get the expected response
	err = json.Unmarshal([]byte(responseRecorder.Body.String()), &js)
	if err != nil {
		t.Error("failed to parse json!")
	}

	// since we specified a start date < 2022-09-09, we expect availability
	if js.Ok || js.Message != "error querying database" {
		t.Error("Got availability when simulating database error")
	}
}

func TestRepository_ReservationSummary(t *testing.T) {
	//reservation data is in session
	resvTest := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "Deluxe suite",
		},
	}

	rq, _ := http.NewRequest("GET", "/reservation-summary", nil)
	ctx := getContext(rq)
	rq = rq.WithContext(ctx)

	responseRecorder := httptest.NewRecorder()
	session.Put(ctx, "reservation", resvTest)

	handler := http.HandlerFunc(Repo.MakeReservationSummary)

	handler.ServeHTTP(responseRecorder, rq)

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("ReservationSummary handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusOK)
	}

	//PostData not in Session
	rq, _ = http.NewRequest("GET", "/reservation-summary", nil)
	ctx = getContext(rq)
	rq = rq.WithContext(ctx)

	responseRecorder = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.MakeReservationSummary)

	handler.ServeHTTP(responseRecorder, rq)

	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("ReservationSummary handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusOK)
	}
}

func TestRepository_SelectAvailableRoom(t *testing.T) {
	//Reservation data  in Session
	resvTest := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "Deluxe suite",
		},
	}

	rq, _ := http.NewRequest("GET", "/select-available-room/1", nil)
	ctx := getContext(rq)
	rq = rq.WithContext(ctx)
	// set the RequestURI on the request so that we can grab the ID
	// from the URL
	rq.RequestURI = "/select-available-room/1"

	responseRecorder := httptest.NewRecorder()
	session.Put(ctx, "reservation", resvTest)

	handler := http.HandlerFunc(Repo.BookRoomNow)

	handler.ServeHTTP(responseRecorder, rq)

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusInternalServerError)
	}

	///*****************************************
	//// second case -- reservation not in session
	//*****************************************/
	rq, _ = http.NewRequest("GET", "/select-available-room/1", nil)
	ctx = getContext(rq)
	rq = rq.WithContext(ctx)
	rq.RequestURI = "/select-available-room/1"

	responseRecorder = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.BookRoomNow)

	handler.ServeHTTP(responseRecorder, rq)

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusOK)
	}

	//Wrong url parameters
	rq, _ = http.NewRequest("GET", "/select-available-room/slim", nil)
	ctx = getContext(rq)
	rq = rq.WithContext(ctx)
	rq.RequestURI = "/select-available-room/slim"

	responseRecorder = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.SelectAvailableRoom)

	handler.ServeHTTP(responseRecorder, rq)

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusOK)
	}
}

func TestRepository_BookRoom(t *testing.T) {
	//Working perfect
	resvTest := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "Deluxe suite",
		},
	}

	rq, _ := http.NewRequest("GET", "/book-room-now?s=2022-09-09&e=2022-09-10&id=1", nil)
	ctx := getContext(rq)
	rq = rq.WithContext(ctx)

	responseRecorder := httptest.NewRecorder()
	session.Put(ctx, "reservation", resvTest)

	handler := http.HandlerFunc(Repo.BookRoomNow)

	handler.ServeHTTP(responseRecorder, rq)

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusOK)
	}

	//Error in the database
	rq, _ = http.NewRequest("GET", "/book-room-now?s=2040-01-01&e=2040-01-02&id=4", nil)
	ctx = getContext(rq)
	rq = rq.WithContext(ctx)

	responseRecorder = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.BookRoomNow)

	handler.ServeHTTP(responseRecorder, rq)

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusInternalServerError)
	}
}

func getContext(rq *http.Request) context.Context {
	ctx, err := session.Load(rq.Context(), rq.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
