package config

import (
	"html/template"
	"log"
	"github.com/alexedwards/scs/v2"
)

//Avoiding creating templates cache all the time a page is display making sure
//it doesn't import anything but can be access in any part of the application
//use in the render & handlers

type AppConfig struct {
	UseCache     bool
	TempCache    map[string]*template.Template
	infoLog      *log.Logger
	InProduction bool
	Session      *scs.SessionManager
}
