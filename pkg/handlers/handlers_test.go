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

	_ "github.com/alexedwards/scs/v2"
	"github.com/dev-ayaa/resvbooking/pkg/driver"
	"github.com/dev-ayaa/resvbooking/pkg/models"
)

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

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	//web server that returns a status code, a test Server and a client must be setup
	sth := httptest.NewTLSServer(routes)

	//close the server after the main function execute
	defer sth.Close()

	for _, h := range hTests {
		//	//this is when the client is requesting to view a page
		//	//h.pagesUrls is the page the client is requesting for
		response, err := sth.Client().Get(sth.URL + h.pagesUrl)
		if err != nil {
			t.Log(err)
			t.Fatalf(fmt.Sprintf("Error Testing the handler of %s......%s", h.pageName, err))
		}
		if response.StatusCode != h.pageStatusCode {
			t.Errorf("Error statusCode for %s : get %v statusCode instead of %v", h.pageName, response.StatusCode, h.pageStatusCode)
		}
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
		correctStatusCode: http.StatusOK,
	},
	{
		testName: "invalid check-in date",
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
		correctStatusCode: http.StatusOK,
	},
	{
		testName: "invalid early reservation",
		postRqData: url.Values{
			"check-in":  {"2000-09-09"},
			"check-out": {"2000-09-10"},
		},
		correctStatusCode: http.StatusOK,
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
		if responseRecorder.Code != m.correctStatusCode {
			t.Errorf("Error Testing for %s in Post availability when no rooms available gave wrong status code: got %d, wanted %d", m.testName, responseRecorder.Code, m.correctStatusCode)
		}
	}
}

var AvailTest = []struct {
	testName      string
	message       string
	postRqData    url.Values
	correctResult bool
}{
	{
		//Valid data but no message
		testName: "available rooms",
		postRqData: url.Values{
			"check-in":  {"2022-09-09"},
			"check-out": {"2022-09-10"},
			"room_id":   {"1"},
		},
		correctResult: true,
	},
	{
		//Valid data but no message
		testName: "not available rooms",
		postRqData: url.Values{
			"check-in":  {"2030-09-09"},
			"check-out": {"2030-09-10"},
			"room_id":   {"1"},
		},
		correctResult: false,
	},
	{
		//No posted data
		testName:      "no data posted",
		message:       "internal server error",
		postRqData:    nil,
		correctResult: false,
	},
	{
		//error from the database
		testName: "Invalid database query",
		message:  "error querying database",
		postRqData: url.Values{
			"check-in":  {"2060-09-09"},
			"check-out": {"2060-09-10"},
			"room_id":   {"1"},
		},
		correctResult: false,
	},
}

func TestRepository_AvailabilityJSON(t *testing.T) {
	//No available rooms
	var rq *http.Request
	for _, m := range AvailTest {
		if m.postRqData != nil {
			rq, _ = http.NewRequest("POST", "/json-availability", strings.NewReader(m.postRqData.Encode()))
		} else {
			rq, _ = http.NewRequest("POST", "/json-availability", nil)
		}
		ctx := getContext(rq)
		rq = rq.WithContext(ctx)

		rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		responseRecorder := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.JsonAvailabilityPage)

		handler.ServeHTTP(responseRecorder, rq)

		// this time we want to parse JSON and get the expected response
		var js ResponseJSON
		err := json.Unmarshal((responseRecorder.Body.Bytes()), &js)
		if err != nil {
			t.Error("failed to parse json!")
		}

		//checking for the correct message
		if js.Ok != m.correctResult {
			t.Errorf("Error for %s  got result to be %v expected %v", m.testName, js.Ok, m.correctResult)
		}
	}
}

var ResvSummTest = []struct {
	testName          string
	correctStatusCode int
	resv              models.Reservation
}{
	{
		testName:          "Reservation summary",
		correctStatusCode: http.StatusOK,
		resv: models.Reservation{
			RoomID: 1,
			Room: models.Room{
				ID:       1,
				RoomName: "Deluxe suite",
			},
		},
	},
	{
		testName:          "Invalid Reservation summary",
		correctStatusCode: http.StatusOK,
		resv:              models.Reservation{},
	},
}

func TestRepository_ReservationSummary(t *testing.T) {
	//reservation data is in session
	for _, m := range ResvSummTest {
		rq, _ := http.NewRequest("GET", "/reservation-summary", nil)
		ctx := getContext(rq)
		rq = rq.WithContext(ctx)
		responseRecorder := httptest.NewRecorder()
		session.Put(ctx, "reservation", m.resv)

		handler := http.HandlerFunc(Repo.MakeReservationSummary)

		handler.ServeHTTP(responseRecorder, rq)

		if responseRecorder.Code != m.correctStatusCode {
			t.Errorf("ReservationSummary handler %s  returned wrong response code: got %d, wanted %d", m.testName, responseRecorder.Code, m.correctStatusCode)
		}
	}

}

var SelectRoomTest = []struct {
	testName          string
	resv              models.Reservation
	id                string
	correctStatusCode int
}{
	{
		testName: "valid Test",
		resv: models.Reservation{
			RoomID: 1,
			Room: models.Room{
				ID:       1,
				RoomName: "Deluxe suite",
			},
		},
		id:                "1",
		correctStatusCode: http.StatusOK,
	},
	{
		testName:          "no reservation",
		id:                "5",
		correctStatusCode: http.StatusOK,
	},
	{
		testName: "invalid details",
		resv: models.Reservation{
			RoomID: 1,
			Room: models.Room{
				ID:       1,
				RoomName: "Deluxe suite",
			},
		},
		id:                "slim",
		correctStatusCode: http.StatusOK,
	},
}

func TestRepository_SelectAvailableRoom(t *testing.T) {
	//Reservation data  in Session
	for _, m := range SelectRoomTest {
		rq, _ := http.NewRequest("GET", fmt.Sprintf("/select-available-room/%s", m.id), nil)
		ctx := getContext(rq)
		rq = rq.WithContext(ctx)
		// set the RequestURI on the request so that we can grab the ID from the URL
		rq.RequestURI = fmt.Sprintf("/select-available-room/%s", m.id)

		responseRecorder := httptest.NewRecorder()
		session.Put(ctx, "reservation", m.resv)

		handler := http.HandlerFunc(Repo.BookRoomNow)

		handler.ServeHTTP(responseRecorder, rq)

		if responseRecorder.Code != m.correctStatusCode {
			t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, m.correctStatusCode)
		}
	}
}

var BookRoomTest = []struct {
	testName          string
	resv              models.Reservation
	correctStatusCode int
	query             string
	postRqData        url.Values
}{
	{
		testName: "valid bookings",
		resv: models.Reservation{
			RoomID: 1,
			Room: models.Room{
				ID:       1,
				RoomName: "Deluxe suite",
			},
		},
		correctStatusCode: http.StatusOK,
		query:             "?s=2022-09-09&e=2022-09-10&id=1",
		postRqData: url.Values{
			"check_in_date":  {"2022-09-09"},
			"check_out_date": {"2022-09-10"},
			"room_id":        {"1"},
		},
	},
	{
		testName: "invalid bookings",
		resv: models.Reservation{
			RoomID: 1,
			Room: models.Room{
				ID:       1,
				RoomName: "Deluxe suite",
			},
		},
		correctStatusCode: http.StatusOK,
		query:             "?s=2044-09-09&e=2044-09-10&id=4",
		postRqData:        nil,
	},
}

func TestRepository_BookRoom(t *testing.T) {
	var rq *http.Request
	for _, m := range BookRoomTest {
		if m.postRqData != nil {
			rq, _ = http.NewRequest("GET", fmt.Sprintf("/book-room-now%s", m.query), strings.NewReader(m.postRqData.Encode()))
		} else {
			rq, _ = http.NewRequest("GET", fmt.Sprintf("/book-room-now%s", m.query), nil)
		}
		ctx := getContext(rq)
		rq = rq.WithContext(ctx)

		responseRecorder := httptest.NewRecorder()
		session.Put(ctx, "reservation", m.resv)

		handler := http.HandlerFunc(Repo.BookRoomNow)

		handler.ServeHTTP(responseRecorder, rq)

		if responseRecorder.Code != m.correctStatusCode {
			t.Errorf("BookRoom handler : %s returned wrong response code: got %d, wanted %d", m.testName, responseRecorder.Code, m.correctStatusCode)
		}
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
	query              string
	correctHTML        string
	correctUrlLocation string
	correctStatusCode  int
}{
	{
		testName:           "valid-testing",
		query:              "?y=2022&m=05",
		correctHTML:        "",
		correctUrlLocation: "",
		correctStatusCode:  http.StatusSeeOther,
	},
	{
		testName:           "invalid-testing",
		query:              "",
		correctHTML:        "",
		correctUrlLocation: "",
		correctStatusCode:  http.StatusSeeOther,
	},
}

func TestRepository_AdminProcessReservation(t *testing.T) {
	for _, s := range ProcessResv {
		rq, _ := http.NewRequest("GET", fmt.Sprintf("/admin/admin-process-reservation/calendar/1/done/%s", s.query), nil)

		ctx := getContext(rq)
		rq = rq.WithContext(ctx)
		// rq.RequestURI = s.correctURL
		responseRecorder := httptest.NewRecorder()

		rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		handler := http.HandlerFunc(Repo.AdminProcessReservation)

		handler.ServeHTTP(responseRecorder, rq)
		if responseRecorder.Code != s.correctStatusCode {
			t.Errorf("Error Testing for process reservation expected %v but get %v", s.correctStatusCode, responseRecorder.Code)
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

func getContext(rq *http.Request) context.Context {
	ctx, err := session.Load(rq.Context(), rq.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
