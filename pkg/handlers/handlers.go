package handlers

import (
	"github.com/Akinleye007/resvbooking/pkg/config"
	"github.com/Akinleye007/resvbooking/pkg/models"
	"github.com/Akinleye007/resvbooking/pkg/render"
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
	rp.App.Session.Put(rq.Context(),"remote_ip", remoteIpAddr)

	render.Template(wr, "home.page.tmpl", &models.TemplateData{})
}


// AboutPage about page  handler
func (rp Repository) AboutPage(wr http.ResponseWriter, rq *http.Request) {
	stringMap := make(map[string]string)

	remoteIpAddr := rp.App.Session.GetString(rq.Context(),"remote_ip")

	stringMap["remote_ip"] = remoteIpAddr
	stringMap["Test"] = "Go Backend developments"

	render.Template(wr, "about.page.tmpl", &models.TemplateData{
		StringData: stringMap,
	})

}
