package forms

type errors map[string][]string

// Add is aReceiver func that Add an error feedback when the
// form field is invalid
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get is a Receiver function that To get an error feedback
//if form field is invalid
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	//return the value in the slices of string
	return es[0]
}
