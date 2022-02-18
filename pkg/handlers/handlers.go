package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Akinleye007/resvbooking/pkg/config"
	"github.com/Akinleye007/resvbooking/pkg/forms"
	"github.com/Akinleye007/resvbooking/pkg/models"
	"github.com/Akinleye007/resvbooking/pkg/render"
	"log"
	"net/http"
)

var Repo *Repository

// Repository struct to store the app
type Repository struct {
	App *config.AppConfig // a struct

}

// NewRepository  create a new repository
func NewRepository(a *config.AppConfig) *Repository {
	return &Repository{App: a}
}

func NewHandlers(r *Repository) {
	Repo = r
}

// HomePage home page handlers & give the handlers a receiver
func (rp *Repository) HomePage(wr http.ResponseWriter, rq *http.Request) {

	remoteIpAddr := rq.RemoteAddr
	rp.App.Session.Put(rq.Context(), "remote_ip", remoteIpAddr)

	render.Template(wr, "home.page.tmpl", &models.TemplateData{}, rq)
}

// AboutPage about page  handler
func (rp Repository) AboutPage(wr http.ResponseWriter, rq *http.Request) {
	render.Template(wr, "about.page.tmpl", &models.TemplateData{}, rq)

}

//ContactPage handler function
func (rp *Repository) ContactPage(wr http.ResponseWriter, rq *http.Request) {

	render.Template(wr, "contact.page.tmpl", &models.TemplateData{}, rq)
}

//JuniorSuitePage  handler function
func (rp *Repository) JuniorSuitePage(wr http.ResponseWriter, rq *http.Request) {

	render.Template(wr, "junior.page.tmpl", &models.TemplateData{}, rq)
}

//PremiumSuitePage handler function
func (rp *Repository) PremiumSuitePage(wr http.ResponseWriter, rq *http.Request) {

	render.Template(wr, "premium.page.tmpl", &models.TemplateData{}, rq)
}

//DeluxeSuitePage handler function
func (rp *Repository) DeluxeSuitePage(wr http.ResponseWriter, rq *http.Request) {

	render.Template(wr, "deluxe.page.tmpl", &models.TemplateData{}, rq)
}

//PenthousePage handler function
func (rp *Repository) PenthousePage(wr http.ResponseWriter, rq *http.Request) {

	render.Template(wr, "penthouse.page.tmpl", &models.TemplateData{}, rq)
}

//ExecutivePage handler function
func (rp *Repository) ExecutivePage(wr http.ResponseWriter, rq *http.Request) {

	render.Template(wr, "executive.page.tmpl", &models.TemplateData{}, rq)
}

//MakeReservationPage handlers function
func (rp *Repository) MakeReservationPage(wr http.ResponseWriter, rq *http.Request) {
	var newReservation models.ReservationData
	data := make(map[string]interface{})
	data["reservationData"] = newReservation
	render.Template(wr, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.NewForm(nil),
	}, rq)
}

func (rp *Repository) PostMakeReservationPage(wr http.ResponseWriter, rq *http.Request) {
	err := rq.ParseForm()
	if err != nil {
		log.Println(err)
	}

	reservationData := models.ReservationData{
		FirstName:   rq.Form.Get("first-name"),
		LastName:    rq.Form.Get("last-name"),
		Email:       rq.Form.Get("Email"),
		PhoneNumber: rq.Form.Get("phoneNumber"),
		Password:    rq.Form.Get("inputPassword4"),
	}

	form := forms.NewForm(rq.PostForm)

	form.HasForm("first-name", rq)
	//form.HasForm("last-name",rq)
	//form.HasForm("phoneNumber",rq)
	//form.HasForm("Email",rq)
	//form.HasForm("inputPassword4",rq)

	if !form.FormValid() {
		data := make(map[string]interface{})
		data["reservationData"] = reservationData
		render.Template(wr, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		}, rq)
		return
	}
}

//CheckAvailabilityPage handler Function
func (rp *Repository) CheckAvailabilityPage(wr http.ResponseWriter, rq *http.Request) {

	render.Template(wr, "check-availability.page.tmpl", &models.TemplateData{}, rq)
}

//PostCheckAvailabilityPage handler function
func (rp *Repository) PostCheckAvailabilityPage(wr http.ResponseWriter, rq *http.Request) {
	//getting the posted value from the form
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

//CheckAvailabilityPage handler Function
func (rp *Repository) JsonAvailabilityPage(wr http.ResponseWriter, rq *http.Request) {

	myResp := ResponseJSON{
		Name:    "Yusuf Akinleye",
		Ok:      true,
		Message: "Available for freelance",
	}

	output, err := json.MarshalIndent(myResp, "", "     ")
	//check for errors
	if err != nil {
		log.Println(err)
	}

	//this type the browser the type of content it is getting
	wr.Header().Set("Content-type", "application/json")
	wr.Write(output)
	//render.Template(wr, "check-availability.page.tmpl", &models.TemplateData{}, rq)
}
