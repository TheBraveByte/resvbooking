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
func (rp Repository) AboutPage(wr http.ResponseWriter, rq *http.Request) {
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

func (rp Repository) LoginPage(wr http.ResponseWriter, rq *http.Request) {
	render.Template(wr, "login.page.tmpl", &models.TemplateData{Form: forms.NewForm(nil)}, rq)

}

//PostLoginPage post user detail in the database
func (rp Repository) PostLoginPage(wr http.ResponseWriter, rq *http.Request) {
	fmt.Println("Logging in user details")
	var email, password string
	//To prevent session fixation during authentication of user login details
	_ = rp.App.Session.RenewToken(rq.Context())
	err := rq.ParseForm()
	if err != nil {
		rp.App.Session.Put(rq.Context(), "errors", "No parsing the login form")
		return
	}

	email = rq.Form.Get("email")
	password = rq.Form.Get("password")
	form := forms.NewForm(rq.PostForm)
	form.Require("email", "password")
	form.ValidLenCharacter("password", 15, rq)
	form.ValidEmail("email")
	if !form.FormValid() {
		render.Template(wr, "login.page.tmpl", &models.TemplateData{Form: forms.NewForm(nil)}, rq)
		return
		//http.Redirect(wr, rq, "/login", http.StatusTemporaryRedirect)
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
	//user, err := rp.DB.GetUserInfoByID(userID)

}

func (rp Repository) LogOutPage(wr http.ResponseWriter, rq *http.Request) {
	rp.App.Session.Destroy(rq.Context())
	rp.App.Session.RenewToken(rq.Context())
	http.Redirect(wr, rq, "/", http.StatusSeeOther)
	return
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
		http.Redirect(wr, rq, "/", http.StatusTemporaryRedirect)
		return
	}

	resv, ok := rp.App.Session.Get(rq.Context(), "reservation").(models.Reservation)
	if !ok {
		rp.App.Session.Put(rq.Context(), "errors", "Error No data for reservation in session")
		http.Redirect(wr, rq, "/", http.StatusTemporaryRedirect)
		return

	}

	roomID, err := strconv.Atoi(rq.Form.Get("room_id"))
	if err != nil {
		rp.App.Session.Put(rq.Context(), "errors", "Error cannot get valid room id")
		http.Redirect(wr, rq, "/", http.StatusTemporaryRedirect)
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
		http.Redirect(wr, rq, "/", http.StatusTemporaryRedirect)
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
		http.Redirect(wr, rq, "/", http.StatusTemporaryRedirect)
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
		http.Redirect(wr, rq, "/", http.StatusTemporaryRedirect)
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
		http.Redirect(wr, rq, "/", http.StatusTemporaryRedirect)
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
		http.Redirect(wr, rq, "/", http.StatusTemporaryRedirect)
		return
	}
	checkOutDate, err := time.Parse(dateLayout, checkOut)
	if err != nil {
		rp.App.Session.Put(rq.Context(), "errors", "cannot parse check-out-date")
		http.Redirect(wr, rq, "/", http.StatusTemporaryRedirect)
		return
	}

	rooms, err := rp.DB.SearchForAvailableRoom(checkInDate, checkOutDate)
	if err != nil {
		rp.App.Session.Put(rq.Context(), "errors", "No available room to reserve")
		http.Redirect(wr, rq, "/", http.StatusTemporaryRedirect)
		return
	}
	for _, room := range rooms {
		rp.App.InfoLog.Println("Rooms Available :: ", room)
	}

	if len(rooms) == 0 {
		rp.App.InfoLog.Println("NO AVAILABLE ROOMS")
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

	//wr.Write([]byte(fmt.Sprintf("Check-in date is %s\nCheck-out date is %s", checkIn, checkOut)))
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

	err := rq.ParseForm()
	if err != nil {
		myResp := ResponseJSON{
			RoomID: "",
			Ok:     false,
			//CheckInDate:  "",
			//CheckOutDate: "",
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
	//if err != nil {
	//	rp.App.Session.Put(rq.Context(), "errors", "Error Parsing the check in date")
	//	return
	//}

	CheckOutDate, _ := time.Parse(layout, cod)
	//if err != nil {
	//	rp.App.Session.Put(rq.Context(), "errors", "Error Parsing the check in date")
	//	return
	//}

	roomID, _ := strconv.Atoi(rq.Form.Get("room_id"))

	isRoomAvailable, err := rp.DB.SearchRoomAvailabileByRoomID(roomID, CheckInDate, CheckOutDate)

	if err != nil {
		myResp := ResponseJSON{
			RoomID: "",
			Ok:     false,
			//CheckInDate:  "",
			//CheckOutDate: "",
			Message: "error querying database",
		}
		output, _ := json.MarshalIndent(myResp, "", "   ")
		//this type the browser the type of content it is getting
		wr.Header().Set("Content-type", "application/json")
		wr.Write(output)
		return
	}
	myResp := ResponseJSON{
		RoomID:       strconv.Itoa(roomID),
		Ok:           isRoomAvailable,
		CheckInDate:  cid,
		CheckOutDate: cod,
		Message:      "",
	}

	//Creating a Json file from struct type
	output, err := json.MarshalIndent(myResp, "", "     ")

	//	myResp := ResponseJSON{
	//		RoomID:       "1",
	//		Ok:           false,
	//		CheckInDate:  "2022-08-09",
	//		CheckOutDate: "2022-08-09",
	//		Message:      "cannot parse form to Json",
	//	}
	//	output, _ := json.MarshalIndent(myResp, "", "   ")
	//	//this type the browser the type of content it is getting
	//	wr.Header().Set("Content-type", "application/json")
	//	wr.Write(output)
	//	return
	//}
	//this type the browser the type of content it is getting
	wr.Header().Set("Content-type", "application/json")
	wr.Write(output)

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
