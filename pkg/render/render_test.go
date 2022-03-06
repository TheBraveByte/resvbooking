package render

import (
	"github.com/dev-ayaa/resvbooking/pkg/models"
	"net/http"
	"testing"
)

var td models.TemplateData

func TestAddDefaultData(t *testing.T) {

	//var td models.TemplateData
	rq, err := getSession()
	if err != nil {
		t.Error("Error Adding Session to the request")
		t.Error(err)
	}
	session.Put(rq.Context(), "flash", "Projects")
	var tadd = AddDefaultData(&td, rq)
	if tadd.Flash != "Projects" {
		t.Error("Testing Default Data function failed")
	}

}

func getSession() (*http.Request, error) {
	//Setting up request with a session
	rq, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		return nil, err
	}
	ctx := rq.Context()
	ctx, _ = session.Load(ctx, rq.Header.Get("X-Session"))
	//check for err or not
	rq = rq.WithContext(ctx)
	return rq, nil

}

func TestTemplate(t *testing.T) {
	templatesPath = "./../../templates"
	var wr myResponse
	//var tmpl string
	tc, err := TemplateCache()
	if err != nil {
		t.Error(err)
	}
	app.TempCache = tc
	rq, err := getSession()
	if err != nil {
		t.Error(err)
		t.Error("Error rendering Templates")
	}

	err = Template(&wr, "deluxe.page.tmpl", &td, rq)
	if err != nil {
		t.Error("Error getting the templates")
		t.Error(err)
	}

}

func TestNewTemplates(t *testing.T) {
	NewTemplates(&testApp)
}
