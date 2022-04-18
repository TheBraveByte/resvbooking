package render

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dev-ayaa/resvbooking/pkg/config"
	"github.com/dev-ayaa/resvbooking/pkg/models"
	"github.com/justinas/nosurf"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var templatesPath = "./templates"
var functions = template.FuncMap{
	//format a dates, currents date
}

/*Storing the templates Cache into the AppConfig struct type, Import the AppConfig as a pointer in the render package back
now use the type store in the AppConfig in the render package ,To keep the stored data updated import the function
where the AppConfig is store in the render package to the main package
*/
var app *config.AppConfig

// NewTemplates  set the config for the templates package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, rq *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(rq)
	td.Warning = app.Session.PopString(rq.Context(), "Warning")
	td.Flash = app.Session.PopString(rq.Context(), "flash")
	td.Error = app.Session.PopString(rq.Context(), "errors")
	//To check if the user id is in the database
	if app.Session.Exists(rq.Context(), "userID") {
		td.IsAuth = 1
	}
	return td
}

// Template RenderTemplates rendering tmpl templates using the cache created
//{response, request, the templates, the data passed to the templates}
func Template(wr http.ResponseWriter, tmpl string, td *models.TemplateData, rq *http.Request) error {
	/*Get the templates cache from the app config from the main.go file ,check if the template exist in the cache
	how to use the right template hold bytes, creating a buffer for  the template and execute
	used to pass default data to all templates */
	tc := map[string]*template.Template{}

	if app.UseCache {
		/* this is used to prevent creating a templates cache of a particular templates all the time when Request for
		 */
		tc = app.TempCache
	} else {
		/*this is first created when our application start running at first
		storing the valid templates, since the App.configure is a pointer, the templates is permanently
		store in it feature know as app.UseCache*/
		tc, _ = TemplateCache()
	}
	//check if the template exist in the cache15
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("cannot create template")
		//return err
	}

	// how to use the right template hold bytes
	buf := new(bytes.Buffer)
	/*creating a buffer for  the template and execute

	td holds the default data we want to pass to a template*/
	td = AddDefaultData(td, rq)
	_ = t.Execute(buf, td)
	// Getting a response
	_, err := buf.WriteTo(wr)
	if err != nil {
		fmt.Println("Error writing Templates to browsers")
		return errors.New("Error Getting the right template")
	}
	return nil
}

// TemplateCache Working with layout and building a template cache
func TemplateCache() (map[string]*template.Template, error) {
	/*go to the templates folder and find all the templates/pages in the folder and layout
	combined all in to a templates that can actually be served
	getting all the templates files use file path except the layout
	*/

	cache := map[string]*template.Template{}
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", templatesPath))
	if err != nil {
		return cache, err
	}

	/*iterating through the pages files paths
	get the base name of the files only
	create a templates set for all the pages except the layout
	check if the templates matches any layout in the templates directory */

	for _, pg := range pages {
		filename := filepath.Base(pg)
		tmp, err := template.New(filename).Funcs(functions).ParseFiles(pg)

		if err != nil {
			return cache, err
		}
		matchTp, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", templatesPath))
		if len(matchTp) > 0 {
			tmp, err = tmp.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", templatesPath))
			if err != nil {
				return cache, err
			}
		}
		cache[filename] = tmp
	}
	return cache, nil

}
