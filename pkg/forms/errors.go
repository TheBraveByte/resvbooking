package forms

type errors map[string][]string

// Set is aReceiver func that Set an error feedback when the
// form field is invalid
func (e errors) Set(formField, message string) {
	e[formField] = append(e[formField], message)
}

// Get is a Receiver function that To get an error feedback
//if form field is invalid
func (e errors) Get(formField string) string {
	// es := e[field]
	if len(e[formField]) == 0 {
		return ""
	}
	//return the value in the slices of string
	return e[formField][0]
}
