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
	"log"
	"net/http"
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
	var newReservation models.Reservation
	data := make(map[string]interface{})
	data["reservationData"] = newReservation
	render.Template(wr, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.NewForm(nil),
		Data: data,
	}, rq)
}

func (rp *Repository) PostMakeReservationPage(wr http.ResponseWriter, rq *http.Request) {
	/*Clients and Server-side Form Validation is process*/
	err := rq.ParseForm()
	if err != nil {
		//log.Println(err)
		helpers.ServerSideError(wr, err)
	}

	reservationData := models.Reservation{
		FirstName:       rq.Form.Get("first-name"),
		LastName:        rq.Form.Get("last-name"),
		Email:           rq.Form.Get("email"),
		PhoneNumber:     rq.Form.Get("phone-number"),
		Password:        rq.Form.Get("inputPassword"),
		ConfirmPassword: rq.Form.Get("inputPassword4"),
	}

	form := forms.NewForm(rq.PostForm)

	form.Require("first-name", "last-name", "phone-number", "email", "inputPassword4", "inputPassword")
	// form.ValidPassword("inputPassword", rq)
	// form.ValidPassword("inputPassword4", rq)
	/*form.HasForm("first_name", rq)
	form.HasForm("last_name", rq)
	form.HasForm("phone-number", rq)
	form.HasForm("email", rq)
	form.HasForm("inputPassword", rq)
	form.HasForm("inputPassword4", rq)
	*/

	form.ValidLenCharacter("first-name", 3, rq)
	form.ValidLenCharacter("last-name", 3, rq)
	form.ValidEmail("email")
	if form.ValidPassword("inputPassword", 10, rq) != form.ValidPassword("inputPassword4", 10, rq) {
		log.Fatal("Incorrect Password....")
		form.Set("inputPassword", "Incorrect Password")
	}

	if !form.FormValid() {
		data := make(map[string]interface{})
		data["reservationData"] = reservationData
		err := render.Template(wr, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		}, rq)
		if err != nil {
			return
		}
		return
	}

	rp.App.Session.Put(rq.Context(), "reservationData", reservationData)
	//redirect the data back to avoid submitting the form more than onece
	http.Redirect(wr, rq, "/make-reservation-data", http.StatusSeeOther)
}

func (rp *Repository) MakeReservationSummary(wr http.ResponseWriter, rq *http.Request) {

	reservationData, ok := rp.App.Session.Get(rq.Context(), "reservationData").(models.Reservation)
	if !ok {
		fmt.Println(ok)
		rp.App.Session.Put(rq.Context(), "error", "session has not reservation")
		http.Redirect(wr, rq, "/", http.StatusTemporaryRedirect)
		log.Println("Error transferring Data")
		return
	}
	data := make(map[string]interface{})
	data["reservationData"] = reservationData

	rp.App.Session.Remove(rq.Context(), "reservationData")
	err := render.Template(wr, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	}, rq)
	if err != nil {
		return
	}

}

//CheckAvailabilityPage handler Function
func (rp *Repository) CheckAvailabilityPage(wr http.ResponseWriter, rq *http.Request) {

	err := render.Template(wr, "check-availability.page.tmpl", &models.TemplateData{}, rq)
	if err != nil {
		return
	}
}

//PostCheckAvailabilityPage handler function
func (rp *Repository) PostCheckAvailabilityPage(wr http.ResponseWriter, rq *http.Request) {
	//getting the posted value from the form with respect to the field
	checkIn := rq.Form.Get("check-in")
	checkOut := rq.Form.Get("check-out")

	wr.Write([]byte(fmt.Sprintf("Check-in date is %s\nCheck-out date is %s", checkIn, checkOut)))
}

//create a json struct interfaces

type ResponseJSON struct {
	Name    string `json:"name"`
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

// JsonAvailabilityPage  handler Function
func (rp *Repository) JsonAvailabilityPage(wr http.ResponseWriter, rq *http.Request) {

	myResp := ResponseJSON{
		Name:    "Yusuf Akinleye",
		Ok:      true,
		Message: "Available for freelance",
	}

	//Creating a Json file from struct type
	output, err := json.MarshalIndent(myResp, "", "     ")
	//check for errors
	if err != nil {
		//log.Println(err)
		helpers.ServerSideError(wr, err)
	}

	//this type the browser the type of content it is getting
	wr.Header().Set("Content-type", "application/json")
	wr.Write(output)
	//render.Template(wr, "check-availability.page.tmpl", &models.TemplateData{}, rq
}
