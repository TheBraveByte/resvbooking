package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/dev-ayaa/resvbooking/pkg/config"
	"github.com/dev-ayaa/resvbooking/pkg/driver"
	"github.com/dev-ayaa/resvbooking/pkg/forms"
	"github.com/dev-ayaa/resvbooking/pkg/helpers"
	"github.com/dev-ayaa/resvbooking/pkg/models"
	"github.com/dev-ayaa/resvbooking/pkg/render"
	"github.com/dev-ayaa/resvbooking/repository"
	"github.com/dev-ayaa/resvbooking/repository/dbRepository"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Repository struct to store the app Config
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepository
}

var Repo *Repository

// NewRepository  create a new repository
func NewRepository(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{App: a,
		DB: dbRepository.NewPostgresRepository(a, db.PSQL)}

}

//
func NewTestRepository(a *config.AppConfig) *Repository {
	return &Repository{App: a,
		DB: dbRepository.NewTestPostgresRepository(a)}

}

func NewHandlers(r *Repository) {
	Repo = r
}

// HomePage home page handlers & give the handlers a receiver
func (rp *Repository) HomePage(wr http.ResponseWriter, rq *http.Request) {
	err := render.Template(wr, "home.page.tmpl", &models.TemplateData{}, rq)
	if err != nil {
		return
	}
}

// AboutPage about page  handler
func (rp *Repository) AboutPage(wr http.ResponseWriter, rq *http.Request) {
	err := render.Template(wr, "about.page.tmpl", &models.TemplateData{}, rq)
	if err != nil {
		return
	}

}

//ContactPage handler function
func (rp *Repository) ContactPage(wr http.ResponseWriter, rq *http.Request) {

	err := render.Template(wr, "contact.page.tmpl", &models.TemplateData{}, rq)
	if err != nil {
		return
	}
}

//CheckAvailabilityPage : This is a Get request handler which only render the check availability page
func (rp *Repository) CheckAvailabilityPage(wr http.ResponseWriter, rq *http.Request) {

	err := render.Template(wr, "check-availability.page.tmpl", &models.TemplateData{}, rq)
	if err != nil {
		return
	}
}

//PostCheckAvailabilityPage : Post request handler for the available date which a certains
// rooms will be available to user to have the room reserve
func (rp *Repository) PostCheckAvailabilityPage(wr http.ResponseWriter, rq *http.Request) {

	err := rq.ParseForm()
	if err != nil {
		rp.App.Session.Put(rq.Context(), "errors", "error cannot parse check availability form")
		http.Redirect(wr, rq, "/", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})

	//getting the posted value from the form with respect to the field
	checkIn := rq.Form.Get("check-in")
	checkOut := rq.Form.Get("check-out")

	//Converting the date in string format to time.Time format
	dateLayout := "2006-01-02"
	checkInDate, err := time.Parse(dateLayout, checkIn)
	if err != nil {
		rp.App.Session.Put(rq.Context(), "errors", "cannot parse check-in-date")
		http.Redirect(wr, rq, "/", http.StatusSeeOther)
		return
	}
	checkOutDate, err := time.Parse(dateLayout, checkOut)
	if err != nil {
		rp.App.Session.Put(rq.Context(), "errors", "cannot parse check-out-date")
		http.Redirect(wr, rq, "/", http.StatusSeeOther)
		return
	}

	rooms, err := rp.DB.SearchForAvailableRoom(checkInDate, checkOutDate)
	if err != nil {
		rp.App.Session.Put(rq.Context(), "errors", "No available room to reserve")
		http.Redirect(wr, rq, "/", http.StatusSeeOther)
		return
	}
	// for _, room := range rooms {
	// 	rp.App.InfoLog.Println("Rooms Available :: ", room)
	// }

	if len(rooms) == 0 {
		// rp.App.InfoLog.Println("NO AVAILABLE ROOMS")
		rp.App.Session.Put(rq.Context(), "errors", "No availale rooms")
		http.Redirect(wr, rq, "/check-availability", http.StatusSeeOther)
		return
	}
	data["rooms"] = rooms

	resv := models.Reservation{
		CheckInDate:  checkInDate,
		CheckOutDate: checkOutDate,
	}

	//After checking for available room by date and store it in session
	rp.App.Session.Put(rq.Context(), "reservation", resv)

	render.Template(wr, "select-available-room.page.tmpl", &models.TemplateData{
		Data: data,
	}, rq)

}

//create a json struct interfaces
type ResponseJSON struct {
	RoomID       string `json:"room_id"`
	Ok           bool   `json:"ok"`
	CheckInDate  string `json:"check_in_date"`
	CheckOutDate string `json:"check_out_date"`
	Message      string `json:"message"`
}

// JsonAvailabilityPage  handler Function
func (rp *Repository) JsonAvailabilityPage(wr http.ResponseWriter, rq *http.Request) {

	//for server error during connection
	err := rq.ParseForm()
	if err != nil {
		myResp := ResponseJSON{
			RoomID:  "",
			Ok:      false,
			Message: "server error",
		}
		output, _ := json.MarshalIndent(myResp, "", "   ")
		//this type the browser the type of content it is getting
		wr.Header().Set("Content-type", "application/json")
		wr.Write(output)
		return

	}
	cid := rq.Form.Get("check-in")
	cod := rq.Form.Get("check-out")

	layout := "2006-01-02"

	CheckInDate, _ := time.Parse(layout, cid)
	CheckOutDate, _ := time.Parse(layout, cod)
	
	roomID, _ := strconv.Atoi(rq.Form.Get("room_id"))

	isRoomAvailable, err := rp.DB.SearchRoomAvailabileByRoomID(roomID, CheckInDate, CheckOutDate)

	//For the database
	myResp := ResponseJSON{
		RoomID:       strconv.Itoa(roomID),
		Ok:           isRoomAvailable,
		CheckInDate:  cid,
		CheckOutDate: cod,
		Message:      "",
	}

	//Creating a Json file from struct type
	output, _ := json.MarshalIndent(myResp, "", "     ")

	//this type the browser the type of content it is getting
	wr.Header().Set("Content-type", "application/json")
	wr.Write(output)

	//For the database
	if err != nil {
		myResp := ResponseJSON{
			RoomID: "",
			Ok:     false,
			Message: "error querying database",
		}
		output, _ := json.MarshalIndent(myResp, "", "   ")
		//this type the browser the type of content it is getting
		wr.Header().Set("Content-type", "application/json")
		wr.Write(output)
		return
	}

	
}

//SelectAvailableRoom : This allow the user to check for all available rooms by RoomID and
//reserve and ony of the available room
func (rp *Repository) SelectAvailableRoom(wr http.ResponseWriter, rq *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(rq, "id"))
	if err != nil {
		rp.App.Session.Put(rq.Context(), "errors", "cannot get room id from the URL")
		return
	}
	resv, ok := rp.App.Session.Get(rq.Context(), "reservation").(models.Reservation)
	if !ok {
		rp.App.Session.Put(rq.Context(), "errors", "cannot get stored data from the reservation database")
		return
	}
	resv.RoomID = roomID
	rp.App.Session.Put(rq.Context(), "reservation", resv)
	http.Redirect(wr, rq, "/make-reservation", http.StatusSeeOther)
}

//BookRoomNow : this handler help the user to book room of their choice right from the
//by checking for availability before signing up to reserve the room choosen
func (rp Repository) BookRoomNow(wr http.ResponseWriter, rq *http.Request) {

	var resv models.Reservation
	room_id, err := strconv.Atoi(rq.URL.Query().Get("room_id"))
	cid := rq.URL.Query().Get("check_in_date")
	cod := rq.URL.Query().Get("check_out_date")

	dateLayout := "2006-01-02"
	checkInDate, err := time.Parse(dateLayout, cod)
	checkOutDate, err := time.Parse(dateLayout, cid)

	//err = errors.New("Error getting data from book-now URL")
	if err != nil {
		rp.App.Session.Put(rq.Context(), "errors", "Error getting data from book-now URL")
		return
	}

	room, err := rp.DB.GetRooms(room_id)
	if err != nil {
		helpers.ServerSideError(wr, err)
		return
	}
	resv.Room.RoomName = room.RoomName
	resv.RoomID = room_id
	resv.CheckInDate = checkInDate
	resv.CheckOutDate = checkOutDate
	rp.App.Session.Put(rq.Context(), "reservation", resv)
	http.Redirect(wr, rq, "/make-reservation", http.StatusSeeOther)

}

//JuniorSuitePage  handler function
func (rp *Repository) JuniorSuitePage(wr http.ResponseWriter, rq *http.Request) {

	err := render.Template(wr, "junior.page.tmpl", &models.TemplateData{}, rq)
	if err != nil {
		return
	}
}

//PremiumSuitePage handler function
func (rp *Repository) PremiumSuitePage(wr http.ResponseWriter, rq *http.Request) {

	err := render.Template(wr, "premium.page.tmpl", &models.TemplateData{}, rq)
	if err != nil {
		return
	}
}

//DeluxeSuitePage handler function
func (rp *Repository) DeluxeSuitePage(wr http.ResponseWriter, rq *http.Request) {

	err := render.Template(wr, "deluxe.page.tmpl", &models.TemplateData{}, rq)
	if err != nil {
		return
	}
}

//PenthousePage handler function
func (rp *Repository) PenthousePage(wr http.ResponseWriter, rq *http.Request) {

	err := render.Template(wr, "penthouse.page.tmpl", &models.TemplateData{}, rq)
	if err != nil {
		return
	}
}

//ExecutivePage handler function
func (rp *Repository) ExecutivePage(wr http.ResponseWriter, rq *http.Request) {

	err := render.Template(wr, "executive.page.tmpl", &models.TemplateData{}, rq)
	if err != nil {
		return
	}
}

//MakeReservationPage handlers function
func (rp *Repository) MakeReservationPage(wr http.ResponseWriter, rq *http.Request) {
	resv, ok := rp.App.Session.Get(rq.Context(), "reservation").(models.Reservation)
	if !ok {
		rp.App.Session.Put(rq.Context(), "errors", "Error linking with Session: No data in session")
		http.Redirect(wr, rq, "/", http.StatusTemporaryRedirect)
		//helpers.ServerSideError(wr, errors.New("error linking with sessions"))
		return
	}

	room, err := rp.DB.GetRooms(resv.RoomID)
	if err != nil {
		rp.App.Session.Put(rq.Context(), "errors", "Error Getting the valid room id")
		http.Redirect(wr, rq, "/", http.StatusTemporaryRedirect)
		return
		//helpers.ServerSideError(wr, err)
		//return
	}
	resv.Room.RoomName = room.RoomName

	data := make(map[string]interface{})
	stringData := make(map[string]string)

	checkInDate := resv.CheckInDate.Format("2006-01-02")
	checkOutDate := resv.CheckOutDate.Format("2006-01-02")

	stringData["check-in"] = checkInDate
	stringData["check-out"] = checkOutDate

	data["reservation"] = resv

	rp.App.Session.Put(rq.Context(), "reservation", resv)

	render.Template(wr, "make-reservation.page.tmpl", &models.TemplateData{
		Form:       forms.NewForm(nil),
		Data:       data,
		StringData: stringData,
	}, rq)
}

//PostMakeReservationPage : Post request information of the user make judicious use of
//Session to pass around data
func (rp *Repository) PostMakeReservationPage(wr http.ResponseWriter, rq *http.Request) {
	err := rq.ParseForm()
	if err != nil {
		rp.App.Session.Put(rq.Context(), "errors", "Error Cannot Parse form data")
		http.Redirect(wr, rq, "/", http.StatusSeeOther)
		return
	}

	resv, ok := rp.App.Session.Get(rq.Context(), "reservation").(models.Reservation)
	if !ok {
		rp.App.Session.Put(rq.Context(), "errors", "Error No data for reservation in session")
		http.Redirect(wr, rq, "/", http.StatusSeeOther)
		return

	}

	roomID, err := strconv.Atoi(rq.Form.Get("room_id"))
	if err != nil {
		rp.App.Session.Put(rq.Context(), "errors", "Error cannot get valid room id")
		http.Redirect(wr, rq, "/", http.StatusSeeOther)
		return
	}

	resv.FirstName = rq.Form.Get("first-name")
	resv.LastName = rq.Form.Get("last-name")
	resv.Email = rq.Form.Get("email")
	resv.PhoneNumber = rq.Form.Get("phone-number")

	form := forms.NewForm(rq.PostForm)

	//Clients and Server-side Form Validation is process*
	form.Require("first-name", "last-name", "phone-number", "email")

	form.ValidLenCharacter("first-name", 3, rq)
	form.ValidLenCharacter("last-name", 3, rq)
	form.ValidEmail("email")

	if !form.FormValid() {
		data := make(map[string]interface{})
		data["reservation"] = resv
		http.Error(wr, "INVALID INPUTS", http.StatusSeeOther)
		err := render.Template(wr, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		}, rq)
		if err != nil {
			return
		}
		return
	}

	NewResvervationID, err := rp.DB.InsertReservation(resv)
	if err != nil {
		rp.App.Session.Put(rq.Context(), "errors", "Error cannot insert reservation in the database")
		http.Redirect(wr, rq, "/", http.StatusSeeOther)
		return
	}
	restriction := models.RoomRestriction{
		ID:            0,
		RoomID:        roomID,
		ReservationID: NewResvervationID,
		RestrictionID: 1,
		CheckInDate:   resv.CheckInDate,
		CheckOutDate:  resv.CheckOutDate,
	}

	err = rp.DB.InsertRoomRestriction(restriction)
	if err != nil {
		rp.App.Session.Put(rq.Context(), "errors", "Cannot insert room restriction in the data")
		http.Redirect(wr, rq, "/", http.StatusSeeOther)
		return
	}

	//Sending mail notification to customer after make a reservation
	mailMsg := models.MailData{
		MailSubject: "Reservation At Rest Tavern Inn",
		Receiver:    resv.Email,
		Sender:      "ayaaakinleye@gmail.com",
		MailContent: `Hi Yusuf <strong>Akinleye</strong>`,
	}

	rp.App.MailChannel <- mailMsg

	notifyOwner := fmt.Sprintf(`<strong>Notification for Reservation at Rest Tavern</strong><br>"+
		"%v %v have secure a reservation of %v from %v to %v`, resv.FirstName, resv.LastName, resv.Room.RoomName,
		resv.CheckInDate.Format("2006-01-02"), resv.CheckOutDate.Format("2006-01-02"))
	mailMsg = models.MailData{
		MailSubject: "Reservation At Rest Tavern Inn",
		Receiver:    "ayaaakinleye@gmail.com",
		Sender:      "ayaaakinleye@gmail.com",
		MailContent: notifyOwner,
	}

	rp.App.MailChannel <- mailMsg
	rp.App.Session.Put(rq.Context(), "reservation", resv)

	//redirect the data back to avoid submitting the form more than once
	http.Redirect(wr, rq, "/make-reservation-data", http.StatusSeeOther)
}

//MakeReservationSummary : Shows all the user information "Fullname, email, Phone Number, check-in-date,
//check-out-date, Room reserved" and many more
func (rp *Repository) MakeReservationSummary(wr http.ResponseWriter, rq *http.Request) {

	resv, ok := rp.App.Session.Get(rq.Context(), "reservation").(models.Reservation)
	if !ok {
		fmt.Println(ok)
		rp.App.Session.Put(rq.Context(), "error", "session has not reservation")
		http.Redirect(wr, rq, "/", http.StatusSeeOther)
		log.Println("Error transferring Data")
		return
	}
	rp.App.Session.Put(rq.Context(), "reservation", resv)

	data := make(map[string]interface{})
	data["reservation"] = resv

	//Remove the stored data in the session
	rp.App.Session.Remove(rq.Context(), "reservation")

	err := render.Template(wr, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	}, rq)
	if err != nil {
		return
	}

}

func (rp *Repository) LoginPage(wr http.ResponseWriter, rq *http.Request) {
	render.Template(wr, "login.page.tmpl", &models.TemplateData{Form: forms.NewForm(nil)}, rq)

}

//PostLoginPage post user detail in the database
func (rp *Repository) PostLoginPage(wr http.ResponseWriter, rq *http.Request) {
	//fmt.Println("Logging in user details")
	var email, password string
	//To prevent session fixation during authentication of user login details
	_ = rp.App.Session.RenewToken(rq.Context())
	err := rq.ParseForm()
	if err != nil {
		rp.App.Session.Put(rq.Context(), "errors", "No parsing the login form")
		return
	}

	//Get the input value from the form and check for authentication
	email = rq.Form.Get("email")
	password = rq.Form.Get("password")
	form := forms.NewForm(rq.PostForm)
	form.Require("email", "password")
	//form.ValidLenCharacter("password", 15, rq)
	form.ValidEmail("email")
	if !form.FormValid() {
		render.Template(wr, "login.page.tmpl", &models.TemplateData{Form: forms.NewForm(nil)}, rq)
		return
	}
	userID, _, err := rp.DB.AuthenticateUser(password, email)
	if err != nil {
		log.Println(err)
		rp.App.Session.Put(rq.Context(), "errors", "log in with valid details")
		http.Redirect(wr, rq, "/login", http.StatusSeeOther)
		return
	}
	rp.App.Session.Put(rq.Context(), "userID", userID)
	rp.App.Session.Put(rq.Context(), "flash", "successfully logged in")
	http.Redirect(wr, rq, "/", http.StatusSeeOther)

}

//LogOutPage this helps to log out the user or admin out of the site
func (rp *Repository) LogOutPage(wr http.ResponseWriter, rq *http.Request) {
	rp.App.Session.Destroy(rq.Context())
	rp.App.Session.RenewToken(rq.Context())
	http.Redirect(wr, rq, "/", http.StatusSeeOther)
	return
}

//AdminPage this is the Administration management page
func (rp *Repository) AdminPage(wr http.ResponseWriter, rq *http.Request) {
	//user := rp.App.Session.Get(rq.Context(),"user")

	err := render.Template(wr, "admin-dashboard.page.tmpl", &models.TemplateData{}, rq)
	if err != nil {
		return
	}
}

//AdminAllReservation this show all the registered user in the administration page
func (rp *Repository) AdminAllReservation(wr http.ResponseWriter, rq *http.Request) {
	allResv, err := rp.DB.AllReservation()
	data := make(map[string]interface{})
	if err != nil {
		helpers.ServerSideError(wr, err)
		//rp.App.Session.Put(rq.Context(), "errors", "no reservation in the database")
		return
	}
	data["reservation"] = allResv
	render.Template(wr, "admin-all-reservation.page.tmpl", &models.TemplateData{
		Data: data,
	}, rq)

}

//AdminNewReservation this shows all the NEW registered user in the Administration page
func (rp *Repository) AdminNewReservation(wr http.ResponseWriter, rq *http.Request) {
	data := make(map[string]interface{})
	newResv, err := rp.DB.AllNewReservation()
	if err != nil {
		helpers.ServerSideError(wr, err)
		return
	}
	data["reservation"] = newResv
	render.Template(wr, "admin-new-reservation.page.tmpl", &models.TemplateData{
		Data: data,
	}, rq)
}

//AdminShowReservation this shows all the reservation information about a particular user
func (rp *Repository) AdminShowReservation(wr http.ResponseWriter, rq *http.Request) {
	StringData := make(map[string]string)
	data := make(map[string]interface{})
	var src string
	var id int

	userInfo := strings.Split(rq.RequestURI, "/")
	id, err := strconv.Atoi(userInfo[len(userInfo)-2])
	if err != nil {
		log.Println("invalid id conversion")
		return
	}

	src = userInfo[len(userInfo)-3]

	month := rq.URL.Query().Get("m")
	year := rq.URL.Query().Get("y")

	userResv, err := rp.DB.ShowUserReservation(id)
	if err != nil {
		helpers.ServerSideError(wr, err)
		return
	}

	data["reservation"] = userResv
	StringData["src"] = src
	StringData["month"] = month
	StringData["year"] = year
	render.Template(wr, "admin-show-reservation.page.tmpl", &models.TemplateData{
		Form:       forms.NewForm(nil),
		Data:       data,
		StringData: StringData,
	}, rq)
}

//PostAdminShowReservation this show the register user info which can be updated
func (rp *Repository) PostAdminShowReservation(wr http.ResponseWriter, rq *http.Request) {
	var src string
	StringData := make(map[string]string)

	err := rq.ParseForm()
	if err != nil {
		rp.App.Session.Put(rq.Context(), "errors", "error parsing forms")
		return
	}
	userInfo := strings.Split(rq.RequestURI, "/")
	//id, err := strconv.Atoi(userInfo[len(userInfo)-1])
	id, err := strconv.Atoi(userInfo[4])

	if err != nil {
		log.Println("invalid id conversion")
		return
	}
	src = userInfo[3]
	StringData["src"] = src
	userResv, err := rp.DB.ShowUserReservation(id)
	if err != nil {
		helpers.ServerSideError(wr, err)
		return
	}

	userResv.FirstName = rq.Form.Get("first_name")
	userResv.LastName = rq.Form.Get("last_name")
	userResv.PhoneNumber = rq.Form.Get("phone_number")
	userResv.Email = rq.Form.Get("email")

	err = rp.DB.UpdateUserReservation(userResv)
	if err != nil {
		rp.App.Session.Put(rq.Context(), "error", "error updating user reservation")
		return
	}
	month := rq.Form.Get("month")
	year := rq.Form.Get("year")
	rp.App.Session.Put(rq.Context(), "flash", "saved")

	if year == "" {
		http.Redirect(wr, rq, fmt.Sprintf("/admin/admin-%s-reservation", src), http.StatusSeeOther)
	} else {
		http.Redirect(wr, rq, fmt.Sprintf("/admin/admin-reservation-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}

}

func (rp Repository) AdminProcessReservation(wr http.ResponseWriter, rq *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(rq, "id"))
	src := chi.URLParam(rq, "src")

	err := rp.DB.ProcessedUpdateReservation(id, 1)
	if err != nil {
		rp.App.Session.Put(rq.Context(), "errors", "unable to processed reservation update")
		return
	}
	year := rq.URL.Query().Get("y")
	month := rq.URL.Query().Get("m")

	rp.App.Session.Put(rq.Context(), "flash", "Reservation marked as processed")

	if year == "" {
		http.Redirect(wr, rq, fmt.Sprintf("/admin/admin-%s-reservations", src), http.StatusSeeOther)
	} else {
		http.Redirect(wr, rq, fmt.Sprintf("/admin/admin-reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}

}

func (rp Repository) AdminDeleteReservation(wr http.ResponseWriter, rq *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(rq, "id"))
	src := chi.URLParam(rq, "src")
	err := rp.DB.DeleteUserReservation(id)
	if err != nil {
		log.Println(err)
		return
	}
	year := rq.URL.Query().Get("y")
	month := rq.URL.Query().Get("m")
	rp.App.Session.Put(rq.Context(), "flash", "Reservation deleted")

	if year == "" {
		http.Redirect(wr, rq, fmt.Sprintf("/admin/admin-%s-reservations", src), http.StatusSeeOther)
	} else {
		http.Redirect(wr, rq, fmt.Sprintf("/admin/admin-reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}
}

//AdminReservationCalendar this shows the calendar schedule of all reservations
func (rp *Repository) AdminReservationCalendar(wr http.ResponseWriter, rq *http.Request) {
	data := make(map[string]interface{})

	present := time.Now()
	if rq.URL.Query().Get("y") != "" {
		year, _ := strconv.Atoi(rq.URL.Query().Get("y"))
		month, _ := strconv.Atoi(rq.URL.Query().Get("m"))
		//func Date(year int, month Month, day, hour, min, sec, nsec int, loc *Location)
		present = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}
	data["present"] = present
	//this increase or add up to the current day // Decrease the date
	nextDate := present.AddDate(0, 1, 0)
	lastDate := present.AddDate(0, -1, 0)

	//Dateformatting
	nextMonthDate := nextDate.Format("01")
	nextMonthYearDate := nextDate.Format("2006")

	lastMonthDate := lastDate.Format("01")
	lastMonthYearDate := lastDate.Format("2006")

	//Storing Formatted date in the database
	StringData := make(map[string]string)
	StringData["next_month_date"] = nextMonthDate
	StringData["next_month_year_date"] = nextMonthYearDate
	StringData["last_month_date"] = lastMonthDate
	StringData["last_month_year_date"] = lastMonthYearDate

	StringData["current_month"] = present.Format("01")
	StringData["current_month_year"] = present.Format("2006")

	//Knowing th current date of the month
	presentYear, presentMonth, _ := present.Date()
	presentLocation := present.Location()
	firstDay := time.Date(presentYear, presentMonth, 1, 0, 0, 0, 0, presentLocation)
	lastDay := firstDay.AddDate(0, 1, -1)

	IntData := make(map[string]int)
	IntData["days_in_month"] = lastDay.Day()

	allRooms, err := rp.DB.AllRoom()
	if err != nil {
		rp.App.Session.Put(rq.Context(), "errors", "error cannot get room schedule from calendar")
		return
	}
	data["rooms"] = allRooms
	for _, room := range allRooms {
		reservationMap := make(map[string]int)
		blockMap := make(map[string]int)

		for d := firstDay; !d.After(lastDay) ; d = d.AddDate(0, 0, 1) {
			reservationMap[d.Format("2006-01-2")] = 0
			blockMap[d.Format("2006-01-2")] = 0
		}

		// get all the restrictions for the current room
		restrictions, err := rp.DB.GetRestrictionsForRoomByDate(room.ID, firstDay, lastDay)
		if err != nil {
			helpers.ServerSideError(wr, err)
			return
		}

		for _, y := range restrictions {
			if y.ReservationID > 0 {
				// it's a reservation with respect to the date
				for d := y.CheckInDate; !d.After(y.CheckOutDate); d = d.AddDate(0, 0, 1) {
					reservationMap[d.Format("2006-01-2")] = y.ReservationID
				}
			} else {
				// it's a block
				blockMap[y.CheckInDate.Format("2006-01-2")] = y.ID
			}
		}
		data[fmt.Sprintf("reservation_map_%d", room.ID)] = reservationMap
		data[fmt.Sprintf("block_map_%d", room.ID)] = blockMap
		fmt.Println(
			StringData["next_month_date"],
			StringData["next_month_year_date"],
			StringData["last_month_date"],
			StringData["last_month_year_date"])

		rp.App.Session.Put(rq.Context(), fmt.Sprintf("block_map_%d", room.ID), blockMap)
	}

	render.Template(wr, "admin-reservation-calendar.page.tmpl", &models.TemplateData{
		StringData: StringData,
		IntData:    IntData,
		Data:       data}, rq)

}

func (rp Repository) PostAdminReservationCalendar(wr http.ResponseWriter, rq *http.Request) {

	err := rq.ParseForm()
	if err != nil {
		helpers.ServerSideError(wr, err)
		return
	}
	year, _ := strconv.Atoi(rq.Form.Get("y"))
	month, _ := strconv.Atoi(rq.Form.Get("m"))

	allRooms, err := rp.DB.AllRoom()
	if err != nil {
		helpers.ServerSideError(wr, err)
		return
	}
	form := forms.NewForm(rq.PostForm)
	//this is to implement the blocked reservations
	for _, r := range allRooms {
		currentMap := rp.App.Session.Get(rq.Context(), fmt.Sprintf("block_map_%d", r.ID)).(map[string]int)
		for key, value := range currentMap {
			if val, ok := currentMap[key]; ok {
				if val > 0 {
					if !form.HasForm(fmt.Sprintf("remove_block_%d_%s", r.ID, key)) {
						err := rp.DB.DeleteBlockByID(value)
						if err != nil {
							log.Println(err)
						}
					}
				}
			}
		}
	} // now handle new blocks
	for name, _ := range rq.PostForm {
		if strings.HasPrefix(name, "add_block") {
			exploded := strings.Split(name, "_")
			roomID, _ := strconv.Atoi(exploded[2])
			t, _ := time.Parse("2006-01-2", exploded[3])
			// insert a new block
			err := rp.DB.InsertBlockForRoom(roomID, t)
			if err != nil {
				log.Println(err)
			}
		}
	}

	rp.App.Session.Put(rq.Context(), "flash", "Changes saved")
	http.Redirect(wr, rq, fmt.Sprintf("/admin/admin-reservation-calendar?y=%d&m=%d", year, month), http.StatusSeeOther)
}
