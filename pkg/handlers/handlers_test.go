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
	{pageName: "AboutPage", pagesUrl: "/about", pageMethod: "GET", pageStatusCode: http.StatusOK},
	{pageName: "ContactPage", pagesUrl: "/contact", pageMethod: "GET", pageStatusCode: http.StatusOK},
	{pageName: "JuniorSuitePage", pagesUrl: "/junior-suite", pageMethod: "GET", pageStatusCode: http.StatusOK},
	{pageName: "DeluxeSuitePage", pagesUrl: "/deluxe-suite", pageMethod: "GET", pageStatusCode: http.StatusOK},
	{pageName: "MakeReservationPage", pagesUrl: "/make-reservation", pageMethod: "GET", pageStatusCode: http.StatusOK},
	{pageName: "MakeReservationSummary", pagesUrl: "/make-reservation-data", pageMethod: "GET", pageStatusCode: http.StatusOK},
	{pageName: "CheckAvailabilityPage", pagesUrl: "/check-availability", pageMethod: "GET", pageStatusCode: http.StatusOK},
	{"LoginPage", "/login", "GET", http.StatusOK},
	{"LogOutPage", "/logout", "GET", http.StatusOK},

	{"Admin", "/admin/dashboard", "GET", http.StatusOK},
	{"NewResvPage", "/admin/admin-new-reservation", "GET", http.StatusOK},
	{"AllResvPage", "/admin/admin-all-reservation", "GET", http.StatusOK},
	{"ResvCalendar", "/admin/admin-reservation-calendar", "GET", http.StatusOK},
	{"ResvCalendarValue", "/admin/admin-reservation-calendar?y=2022&m=05", "GET", http.StatusOK},

	{"ShowResv", "/admin/admin-show-reservation/new/1/show", "GET", http.StatusOK},
	//{"DeleteResv", "/admin/admin-delete-reservation/new/1/done", "GET", http.StatusSeeOther},
	//{"ProcessResv", "/admin/admin-process-reservation/new/1/done", "GET", http.StatusSeeOther},

	//{"PostShowResv", "/admin/admin-show-reservation/new/1", "POST", http.StatusSeeOther},
	//{"PostLoginPage", "/login", "POST", http.StatusSeeOther},
	//{"PostResvCalendar", "/admin/admin-reservation-calendar", "POST", http.StatusSeeOther},
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

var MakeResvTest = []struct {
	resv              models.Reservation
	correctStatusCode int
}{
	{
		models.Reservation{
			RoomID: 1,
			Room: models.Room{
				ID:       1,
				RoomName: "Deluxe suite",
			},
		},
		http.StatusOK,
	},
}

func TestRepository_MakeReservationPage(t *testing.T) {

	for _, m := range MakeResvTest {
		rq, _ := http.NewRequest("GET", "/make-reservation", nil)
		ctx := getContext(rq)
		rq = rq.WithContext(ctx)

		responseRecorder := httptest.NewRecorder()
		session.Put(ctx, "reservation", m.resv)
		handler := http.HandlerFunc(Repo.MakeReservationPage)
		handler.ServeHTTP(responseRecorder, rq)

		if responseRecorder.Code != m.correctStatusCode {
			t.Errorf("Wrong response from the make-reservation handler: got %v wanted %v", responseRecorder.Code, m.correctStatusCode)
		}
	}

}

var PostMakeResv = []struct {
	testName          string
	postRqData        url.Values
	correctStatusCode int
}{
	{
		testName: "valid reservation",
		postRqData: url.Values{
			"check-in":     {"2022-03-01"},
			"check-out":    {"2022-03-06"},
			"first-name":   {"Graham"},
			"last-name":    {"Graham"},
			"email":        {"Grahams@gmail.com"},
			"phone-number": {"20229028844"},
			"room_id":      {"1"},
		},
		correctStatusCode: http.StatusSeeOther,
	},

	{
		testName: "invalid reservation check-in date",
		postRqData: url.Values{
			"check-in":     {"invalid"},
			"check-out":    {"2022-09-10"},
			"first-name":   {"Graham"},
			"last-name":    {"Graham"},
			"email":        {"Grahams@gmail.com"},
			"phone-number": {"20229028844"},
			"room_id":      {"1"},
		},
		correctStatusCode: http.StatusSeeOther,
	},
	{
		testName: "invalid reservation check-out date",
		postRqData: url.Values{
			"check-in":     {"2022-09-09"},
			"check-out":    {"invalid"},
			"first-name":   {"Graham"},
			"last-name":    {"Graham"},
			"email":        {"Grahams@gmail.com"},
			"phone-number": {"20229028844"},
			"room_id":      {"1"},
		},
		correctStatusCode: http.StatusSeeOther,
	},
	{
		testName: "invalid reservation room id",
		postRqData: url.Values{
			"check-in":     {"2022-09-09"},
			"check-out":    {"2022-09-10"},
			"first-name":   {"Graham"},
			"last-name":    {"Graham"},
			"email":        {"Grahams@gmail.com"},
			"phone-number": {"20229028844"},
			"room_id":      {"invalid"},
		},
		correctStatusCode: http.StatusSeeOther,
	},
	{
		testName: "invalid reservation data",
		postRqData: url.Values{
			"check-in":     {"2022-09-09"},
			"check-out":    {"2022-09-10"},
			"first-name":   {"G"},
			"last-name":    {"G"},
			"email":        {"Grahams@gmail.com"},
			"phone-number": {"20229028844"},
			"room_id":      {"1"},
		},
		correctStatusCode: http.StatusSeeOther,
	},
	{
		testName: "insert invalid reservation",
		postRqData: url.Values{
			"check-in":     {"2022-09-09"},
			"check-out":    {"2022-09-10"},
			"first-name":   {"Graham"},
			"last-name":    {"Graham"},
			"email":        {"Grahams@gmail.com"},
			"phone-number": {"20229028844"},
			"room_id":      {"14"},
		},
		correctStatusCode: http.StatusSeeOther,
	},
	{
		testName: "insert room restriction",
		postRqData: url.Values{
			"check-in":     {"invalid"},
			"check-out":    {"2022-09-10"},
			"first-name":   {"Graham"},
			"last-name":    {"Graham"},
			"email":        {"Grahams@gmail.com"},
			"phone-number": {"20229028844"},
			"room_id":      {"11"},
		},
		correctStatusCode: http.StatusSeeOther,
	},
}

func TestRepository_PostMakeReservationPage(t *testing.T) {
	for _, m := range PostMakeResv {
		rq, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(m.postRqData.Encode()))
		ctx := getContext(rq)
		rq = rq.WithContext(ctx)
		rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		responseRecorder := httptest.NewRecorder()
		handler := http.HandlerFunc(Repo.PostMakeReservationPage)
		handler.ServeHTTP(responseRecorder, rq)

		if responseRecorder.Code != m.correctStatusCode {
			t.Errorf("Wrong response for %s from the make-reservation handler: got %v wanted %v", m.testName, responseRecorder.Code, m.correctStatusCode)
		}
	}
}
func TestNewRepo(t *testing.T) {
	var db driver.DB
	testNewRepo := NewRepository(&app, &db)
	if reflect.TypeOf(testNewRepo).String() != "*handlers.Repository" {
		t.Errorf("Did not get the correct repository : got %v wanted *Repository", reflect.TypeOf(NewTestRepository).String())
	}
}

var CheckAvailTest = []struct {
	testName          string
	postRqData        url.Values
	correctStatusCode int
}{

	{
		testName: "valid reservation date",
		postRqData: url.Values{
			"check-in":  {"2022-09-09"},
			"check-out": {"2022-09-10"},
		},
		correctStatusCode: http.StatusSeeOther,
	},
	{
		testName: "valid check-in date",
		postRqData: url.Values{
			"check-in":  {"invalid"},
			"check-out": {"2022-09-10"},
		},
		correctStatusCode: http.StatusSeeOther,
	},
	{
		testName: "valid check-out date",
		postRqData: url.Values{
			"check-in":  {"2022-09-09"},
			"check-out": {"2022-09-10"},
		},
		correctStatusCode: http.StatusSeeOther,
	},
	{
		testName: "invalid early reservation",
		postRqData: url.Values{
			"check-in":  {"2000-09-09"},
			"check-out": {"2000-09-10"},
		},
		correctStatusCode: http.StatusSeeOther,
	},
	{
		testName: "invalid future reservation",
		postRqData: url.Values{
			"check-in":  {"2029-09-09"},
			"check-out": {"2029-09-10"},
		},
		correctStatusCode: http.StatusSeeOther,
	},
}

func TestRepository_PostCheckAvailabilityPage(t *testing.T) {
	for _, m := range CheckAvailTest {
		rq, _ := http.NewRequest("POST", "/check-availability", strings.NewReader(m.postRqData.Encode()))
		ctx := getContext(rq)
		rq = rq.WithContext(ctx)
		rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		responseRecorder := httptest.NewRecorder()
		handler := http.HandlerFunc(Repo.PostCheckAvailabilityPage)
		handler.ServeHTTP(responseRecorder, rq)
		if responseRecorder.Code != http.StatusOK {
			t.Errorf("Error Testing for %s in Post availability when no rooms available gave wrong status code: got %d, wanted %d",m.testName, responseRecorder.Code,m.correctStatusCode)
		}
	}
	// PostRqData := url.Values{}

	// PostRqData.Add("check-in", "2022-09-09")
	// PostRqData.Add("check-out", "2022-09-10")

	// //  rooms are available that does not exist
	// PostRqData.Add("check-in", "2000-09-09")
	// PostRqData.Add("check-out", "2000-09-10")

	// rq, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(PostRqData.Encode()))

	// ctx = getContext(rq)
	// rq = rq.WithContext(ctx)

	// rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// responseRecorder = httptest.NewRecorder()
	// handler = http.HandlerFunc(Repo.PostCheckAvailabilityPage)
	// handler.ServeHTTP(responseRecorder, rq)

	// // since we have rooms available, we expect to get status http.StatusOK
	// if responseRecorder.Code != http.StatusOK {
	// 	t.Errorf("Post availability when rooms are available gave wrong status code: got %d, wanted %d", responseRecorder.Code, http.StatusOK)
	// }

	// //No date to check for reservation
	// rq, _ = http.NewRequest("POST", "/search-availability", nil)

	// ctx = getContext(rq)
	// rq = rq.WithContext(ctx)

	// rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// responseRecorder = httptest.NewRecorder()

	// handler = http.HandlerFunc(Repo.PostCheckAvailabilityPage)

	// handler.ServeHTTP(responseRecorder, rq)

	// // since we have rooms available, we expect to get status http.StatusTemporaryRedirect
	// if responseRecorder.Code != http.StatusTemporaryRedirect {
	// 	t.Errorf("Post availability with empty request body (nil) gave wrong status code: got %d, wanted %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	// }

	// //Wrong  Date format
	// // start date in the wrong format
	// PostRqData.Add("check-in", "invalid")
	// PostRqData.Add("check-out", "2022-09-09")

	// rq, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(PostRqData.Encode()))

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

	// // since we have rooms available, we expect to get status http.StatusTemporaryRedirect
	// if responseRecorder.Code != http.StatusOK {
	// 	t.Errorf("Post availability with invalid start date gave wrong status code: got %d, wanted %d", responseRecorder.Code, http.StatusOK)
	// }

	// //Wrong date format for check-out date
	// PostRqData.Add("check-in", "2022-09-09")
	// PostRqData.Add("check-out", "invalid")
	// rq, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(PostRqData.Encode()))

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

	// // since we have rooms available, we expect to get status http.StatusTemporaryRedirect
	// if responseRecorder.Code != http.StatusOK {
	// 	t.Errorf("Post availability with invalid end date gave wrong status code: got %d, wanted %d", responseRecorder.Code, http.StatusOK)
	// }

	// // database Error fails
	// PostRqData.Add("check-in", "2029-09-09")
	// PostRqData.Add("check-out", "2029-09-10")
	// rq, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(PostRqData.Encode()))

	// ctx = getContext(rq)
	// rq = rq.WithContext(ctx)
	// rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// responseRecorder = httptest.NewRecorder()
	// handler = http.HandlerFunc(Repo.PostCheckAvailabilityPage)

	// handler.ServeHTTP(responseRecorder, rq)

	// // since we have rooms available, we expect to get status http.StatusTemporaryRedirect
	// if responseRecorder.Code != http.StatusOK {
	// 	t.Errorf("Post availability when database query fails gave wrong status code: got %d, wanted %d", responseRecorder.Code, http.StatusOK)
	// }
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

//For login test
var loginTests = []struct {
	testName           string
	email              string
	correctStatusCode  int
	correctHTML        string
	correctUrlLocation string
}{
	{
		"valid-login-details",
		"dev-ayaa007@admin.com",
		http.StatusSeeOther,
		"",
		"/",
	},
	{
		"invalid-login-details",
		"dev-ayaa007@jingle.com",
		http.StatusSeeOther,
		"",
		"/login",
	},
	{
		"invalid-data",
		"dev-ayaa007",
		http.StatusOK,
		`action="/login"`,
		"",
	},
}

func TestRepository_PostLoginPage(t *testing.T) {
	for _, d := range loginTests {
		postRqData := url.Values{}
		postRqData.Add("email", d.email)
		postRqData.Add("password", "2701Akin1234")

		rq, _ := http.NewRequest("POST", "/login", strings.NewReader(postRqData.Encode()))
		ctx := getContext(rq)
		rq = rq.WithContext(ctx)
		rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		responseRecorder := httptest.NewRecorder()
		handler := http.HandlerFunc(Repo.PostLoginPage)
		handler.ServeHTTP(responseRecorder, rq)

		//checking for the correct statusCode
		if responseRecorder.Code != d.correctStatusCode {
			t.Errorf("Error Testing LoginPage got wrong ressponse code expected %v got %v", d.correctStatusCode, responseRecorder.Code)
		}

		//checking for the correct html
		if d.correctHTML != "" {
			html := responseRecorder.Body.String()
			if !strings.Contains(html, d.correctHTML) {
				t.Errorf("Error invalid templates for login-page expected %v", d.correctHTML)
			}

		}
		//checking for the correct location to be redirect
		if d.correctUrlLocation != "" {
			urlLocation, _ := responseRecorder.Result().Location()
			if urlLocation.String() != d.correctUrlLocation {
				t.Errorf("Error invalid location for login-page expected %v", d.correctUrlLocation)
			}
		}

	}
}

var ShowResvTest = []struct {
	testName           string
	correctURL         string
	correctUrlLocation string
	correctStatusCode  int
	correctHTML        string
	postRqData         url.Values
}{
	{
		"valid-reservation",
		"/admin/admin-all-reservation/all/1/show",
		"/admin/admin-all-reservation",
		http.StatusSeeOther,
		"",
		url.Values{
			"first_name":   {"Yusuf"},
			"last_name":    {"Akinleye"},
			"phone_number": {"09088765312"},
			"email":        {"dev-ayaa007@admin.com"},
		},
	},
	//{
	//	"valid-reservation",
	//	"/admin/admin-reservation-calendar/calendar/1",
	//	"/admin/admin-reservation-calendar/calendar?y=2022&m=05",
	//	http.StatusSeeOther,
	//	"",
	//	url.Values{
	//		"first_name":   {"Yusuf"},
	//		"last_name":    {"Akinleye"},
	//		"phone_number": {"09088765312"},
	//		"email":        {"dev-ayaa007@admin.com"},
	//		"month":        {"05"},
	//		"year":         {"2022"},
	//	},
	//},
	{
		"valid-reservation",
		"/admin/admin-new-reservation/new/1/show",
		"/admin/admin-new-reservation",
		http.StatusSeeOther,
		"",
		url.Values{
			"first_name":   {"Yusuf"},
			"last_name":    {"Akinleye"},
			"phone_number": {"09088765312"},
			"email":        {"dev-ayaa007@admin.com"},
		},
	},
}

func TestRepository_PostAdminShowReservation(t *testing.T) {
	for _, s := range ShowResvTest {
		var rq *http.Request
		if s.postRqData != nil {
			rq, _ = http.NewRequest("POST", "/login", strings.NewReader(s.postRqData.Encode()))
		} else {
			rq, _ = http.NewRequest("POST", "/login", nil)

		}

		ctx := getContext(rq)
		rq = rq.WithContext(ctx)
		rq.RequestURI = s.correctURL
		responseRecorder := httptest.NewRecorder()

		rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		handler := http.HandlerFunc(Repo.PostAdminShowReservation)

		handler.ServeHTTP(responseRecorder, rq)

		if responseRecorder.Code != s.correctStatusCode {
			t.Errorf("Error Testing for postAdminShowReservation expected %v but get %v", s.correctStatusCode, responseRecorder.Code)
		}
		if s.correctUrlLocation != "" {
			urlLocation, _ := responseRecorder.Result().Location()
			if urlLocation.String() != s.correctUrlLocation {
				t.Errorf("Error Testing for invalid url")
			}
		}

		if s.correctHTML != "" {
			html := responseRecorder.Body.String()
			if !strings.Contains(html, s.correctHTML) {
				t.Errorf("Error Testing for invalid html template")
			}
		}

	}
}

var ProcessResv = []struct {
	testName           string
	correctURL         string
	correctHTML        string
	correctUrlLocation string
	postRqData         url.Values
	correctStatusCode  int
}{
	{
		"valid-testing",
		"",
		"",
		"",
		url.Values{
			"month": {"05"},
			"year":  {"2022"},
		},
		http.StatusSeeOther,
	},

	//{
	//	"valid-testing",
	//	"/admin/admin-new-reservation/new/1/show",
	//	"",
	//	"/admin/admin-new-reservation",
	//	url.Values{
	//		"month": {"05"},
	//		"year":  {"2022"},
	//	},
	//	http.StatusSeeOther,
	//},
	//{
	//	"valid-testing",
	//	"/admin/admin-all-reservation/all/1/show",
	//	"",
	//	"/admin/admin-all-reservation",
	//	url.Values{
	//		"month": {"05"},
	//		"year":  {"2022"},
	//	},
	//	http.StatusSeeOther,
	//},
	{
		"valid-testing",
		"?y=2022&m=05",
		"",
		"",
		url.Values{
			"month": {"05"},
			"year":  {"2022"},
		},
		http.StatusSeeOther,
	},
}

func TestRepository_AdminProcessReservation(t *testing.T) {
	for _, s := range ProcessResv {
		var rq *http.Request
		if s.postRqData != nil {
			rq, _ = http.NewRequest("POST", fmt.Sprintf("/admin/admin-process-reservation/calendar/1/done%s", s.correctURL), strings.NewReader(s.postRqData.Encode()))
		} else {
			rq, _ = http.NewRequest("POST", fmt.Sprintf("/admin/admin-process-reservation/calendar/1/done%s", s.correctURL), nil)

		}

		ctx := getContext(rq)
		rq = rq.WithContext(ctx)
		//rq.RequestURI = s.correctURL
		responseRecorder := httptest.NewRecorder()

		rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		handler := http.HandlerFunc(Repo.PostAdminShowReservation)

		handler.ServeHTTP(responseRecorder, rq)
		if responseRecorder.Code != s.correctStatusCode {
			t.Errorf("Error Testing for postAdminShowReservation expected %v but get %v", s.correctStatusCode, responseRecorder.Code)
		}
		//if s.correctUrlLocation != "" {
		//	urlLocation, _ := responseRecorder.Result().Location()
		//	if urlLocation.String() != s.correctUrlLocation {
		//		t.Errorf("Error Testing for invalid url")
		//	}
		//}
		//
		//if s.correctHTML != "" {
		//	html := responseRecorder.Body.String()
		//	if !strings.Contains(html, s.correctHTML) {
		//		t.Errorf("Error Testing for invalid html template")
		//	}
		//}

	}

}
func getContext(rq *http.Request) context.Context {
	ctx, err := session.Load(rq.Context(), rq.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
