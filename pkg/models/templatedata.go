package models

import "github.com/Akinleye007/resvbooking/pkg/forms"

//TemplateData  holds data sent from handlers to templates
type TemplateData struct {
	StringData map[string]string
	intData    map[string]int
	floatData  map[string]float64
	Data       map[string]interface{}
	CSRFToken  string
	popMessage string
	Warning    string
	Flash      string
	Error      string
	Form       *forms.Form
}
