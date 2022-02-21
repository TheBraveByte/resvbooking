package forms

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/http"
	"net/url"
	"strings"
)

type Form struct {
	url.Values
	Error errors
}

// NewForm Initialize a form of type struct
func NewForm(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}

}

//HasForm to check if the form is not empty
func (f *Form) HasForm(field string, rq *http.Request) bool {
	/*This to check if the form input has value or not*/
	checkForm := rq.Form.Get(field)
	if checkForm == "" {
		f.Error.Add(field, "this field cannot be blank")
		return false
	}
	return true

}

// FormValid Validate the forms values
func (f *Form) FormValid() bool {

	return len(f.Error) == 0
}

func (f *Form) Require(field ...string) {
	for _, afield := range field {
		value := f.Get(afield)
		if strings.TrimSpace(value) == "" {
			f.Error.Add(afield, "This field cant be blank")
		}
	}
}

// ValidLenCharacter This check for valid minimum length of input character in the form field
func (f *Form) ValidLenCharacter(field string, CharLen int, rq *http.Request) bool {
	fd := rq.Form.Get(field)
	if len(fd) < CharLen {
		f.Error.Add(field, fmt.Sprintf("This field must have at least %d character long", CharLen))
		return false
	}
	return true
}

// ValidEmail check for valid email
func (f *Form) ValidEmail(field string) bool {
	checkEMail := f.Get(field)
	//use Govalidator for email
	if !govalidator.IsEmail(checkEMail) {
		f.Error.Add(field, "Invalid Email Address")
		return false
	}
	return true
}
