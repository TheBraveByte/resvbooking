package render

import (
	"bytes"
	"fmt"
	"github.com/Akinleye007/resvbooking/pkg/config"
	"github.com/Akinleye007/resvbooking/pkg/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{
	//format a dates, currents date
}

//Storing the templates Cache into the AppConfig struct type
//Import the AppConfig as a pointer in the render package back
//now use the type store in the AppConfig in the render package
//To keep the stored data updated import the function where the AppConfig is store in the render package to the main package
var app *config.AppConfig

// NewTemplates  set the config for the templates package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// Template RenderTemplates rendering tmpl templates using the cache created
func Template(wr http.ResponseWriter, tmpl string, td *models.TemplateData) {
	//Get the templates cache from the app config from the main.go file
	tc := map[string]*template.Template{}

	if app.UseCache {
		tc = app.TempCache
	} else {
		tc, _ = TemplateCache()
	} //tc:= app.TempCache // tc is a map
	//tc, err := TemplateCache()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//check if the template exist in the cache15
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("cannot create template")
	}

	// how to use the right template hold bytes
	buf := new(bytes.Buffer)

	// creating a buffer for  the template and execute
	td = AddDefaultData(td) // use to pass default data to all templates


	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(wr)
	if err !=
		nil {
		fmt.Println("Error writing Templates to browsers")
	}

	//parseTemplate, _ := template.ParseFiles("./templates/" + tmpl)
	//err = parseTemplate.Execute(wr, nil)
	//if err != nil {
	//	fmt.Println("Error Rendering Templates")
	//	return
	//}

}

// TemplateCache Working with layout and building a template cache
func TemplateCache() (map[string]*template.Template, error) {

	//go to the templates folder and find all the templates/pages in the folder and layout
	//combined all in to a templates that can actually be served

	cache := map[string]*template.Template{}

	//getting all the templates files use file path except the layout
	pages, err := filepath.Glob("./templates/*.page.tmpl")

	if err != nil {
		return cache, err
	}

	//iterating through the pages files paths
	for _, pg := range pages {
		fmt.Println("Pages",pg)

		//get the base name of the files only
		filename := filepath.Base(pg)
		fmt.Printf("Loading %v currently\n", filename)
		fmt.Println("Cache",cache)

		//create a templates set for all the pages except the layout
		tpset, err := template.New(filename).Funcs(functions).ParseFiles(pg)
		fmt.Println("templates set : ",tpset)

		if err != nil {
			return cache, err
		}

		// check if the templates matches any layout in the templates directory
		matchTp, err := filepath.Glob("./templates/*.layout.tmpl")
		fmt.Println(matchTp)

		if len(matchTp) > 0 {
			tpset, err = tpset.ParseGlob("./templates/*.layout.tmpl")
			fmt.Println("Template ---->",tpset)

			if err != nil {
				return cache, err
			}
		}
		cache[filename] = tpset

	}
	return cache, nil

}
