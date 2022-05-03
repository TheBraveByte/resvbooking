package models

import "github.com/dev-ayaa/resvbooking/pkg/forms"

//TemplateData  holds data sent from handlers to templates
type TemplateData struct {
	StringData map[string]string
	IntData    map[string]int
	FloatData  map[string]float64
	Data       map[string]interface{}
	CSRFToken  string
	PopMessage string
	Warning    string
	Flash      string
	Error      string
	Form       *forms.Form
	IsAuth     int
}
