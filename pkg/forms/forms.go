package forms

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/http"
	"net/url"
	"strings"
)

type Form struct {
	url.Values // data typed in the form
	Error      errors
}

// NewForm Initialize a form of type struct
func NewForm(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) Require(formField ...string) {
	for _, field := range formField {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Error.Set(field, "This field can't be blank")
		}
	}
}

//HasForm to check if the form is not empty
func (f *Form) HasForm(formField string) bool {
	/*This to check if the form input has value or not*/
	checkForm := f.Get(formField)
	if checkForm == "" {
		f.Error.Set(formField, "this field cannot be blank")
		return false
	}
	return true

}

// FormValid Validate the forms values
func (f *Form) FormValid() bool {

	return len(f.Error) == 0
}

// ValidLenCharacter This check for valid minimum length of input character in the form field
func (f *Form) ValidLenCharacter(formField string, CharLen int, rq *http.Request) bool {
	fd := rq.Form.Get(formField)
	if len(fd) < CharLen {
		f.Error.Set(formField, fmt.Sprintf("This field must have at least %d character long", CharLen))
		return false
	}
	return true
}

//ValidEmail check for valid email
func (f *Form) ValidEmail(formField string) bool {
	checkEMail := f.Get(formField)
	//use Govalidator for email
	if !govalidator.IsEmail(checkEMail) {
		f.Error.Set(formField, "Invalid Email Address")
		return false
	}
	return true
}

// func (f *Form) ValidPassword(formField1 string,formField2 string ,minLen int,rq *http.Request) bool {
// 		checkPassword   := f.Get(formField1)
// 		confirmPassword := f.Get(formField2)
// 		if len(checkPassword) != minLen {
// 			f.Error.Set(formField1,"Weak Password : Enter a Minimium of 10 character")
// 			return false
// 		}
// 		if len(confirmPassword)  != minLen {
// 			f.Error.Set(formField2,"Weak Password : Enter a Minimium of 10 character")
// 			return false
// 		}
// 		if checkPassword  != confirmPassword{
// 			f.Error.Set(formField2,"Invalid Password... Enter the correct password")
// 			return false
// 		}
// 		return true
// }

func (f *Form) ValidPassword(formField string, minLen int, rq *http.Request) bool {
	CheckPassword := f.Get(formField)
	if len(CheckPassword) < minLen {
		f.Error.Set(formField, "Weak Password : Enter a Minimium of 10 character")
		return false
	}
	return true

}
