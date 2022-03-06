package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

//var formField string
var postForm url.Values

func TestNewForm(t *testing.T) {
	rq := httptest.NewRequest("POST", "/", nil)
	NewForm(rq.PostForm)

}

func TestForm_FormValid(t *testing.T) {
	rq := httptest.NewRequest("POST", "/", nil)
	form := NewForm(rq.PostForm)
	validForm := form.FormValid()
	if !validForm {
		t.Error("Error Validating the POST form details@Form_Valid")
	}

}

func TestForm_HasForm(t *testing.T) {
	//formField= "Akinleye"
	rq := httptest.NewRequest("POST", "/", nil)
	form := NewForm(rq.PostForm)
	hasForm := form.HasForm("firstname", rq)
	if hasForm {
		t.Error("Error Check for Entry form data @Has_form testing")
	}
	postForm := url.Values{}
	postForm.Add("firstname", "Yusuf")
	//rsp, err := http.Get("firstname")
	rq, _ = http.NewRequest("POST", "/", nil)
	rq.PostForm = postForm
	form = NewForm(rq.PostForm)
	hasForm = form.FormValid()
	if !hasForm {
		t.Error("Error not input data in the form @Has_Form Testing")
	}

}

func TestForm_Require(t *testing.T) {
	rq := httptest.NewRequest("POST", "/", nil)
	form := NewForm(rq.PostForm)
	form.Require("firstname", "lastname", "email")

	if form.FormValid() {
		t.Error("Error the form field is required @form_Required Testing")
	}
	postForms := url.Values{}

	postForms.Add("firstname", "Yusuf")
	postForms.Add("lastname", "Akinleye")
	postForms.Add("email", "Aaa@gmail.com")
	rq, _ = http.NewRequest("POST", "/", nil)
	rq.PostForm = postForms
	form = NewForm(rq.PostForm)
	form.Require("firstname", "lastname", "email")

	if !form.FormValid() {
		t.Error("Error In the Form @Form_Required....Testing")
	}
}

func TestForm_ValidEmail(t *testing.T) {
	rq := httptest.NewRequest("POST", "/", nil)
	form := NewForm(rq.PostForm)
	isValidEmail := form.ValidEmail("email")
	if isValidEmail && form.FormValid() {
		t.Error("Error (No Value) Validating the input email @Valid_Email Testing")
	}

	postForm = url.Values{}
	postForm.Add("email", "YusufAkinleyeAbidemi@gmail.com")

	rq, _ = http.NewRequest("POST", "/", nil)
	//ignore the error
	rq.PostForm = postForm
	form = NewForm(rq.PostForm)
	isValidEmail = form.ValidEmail("email")
	if !form.FormValid() && !isValidEmail {
		t.Error("Error Invalid Email @ Valid_Email Testing")
	}

}

func TestForm_ValidLenCharacter(t *testing.T) {
	rq := httptest.NewRequest("POST", "/", nil)
	form := NewForm(rq.PostForm)
	validChar := form.ValidLenCharacter("firstname", 5, rq)
	if validChar && form.FormValid() {
		t.Error("Error No data in the form field yet @Valid_Len_Character")
	}
	rq, _ = http.NewRequest("POST", "/", nil)
	postForm = url.Values{}
	postForm.Add("firstname", "Yusuf")
	form = NewForm(rq.PostForm)
	rq.PostForm = postForm
	validChar = form.ValidLenCharacter("firstname", 5, rq)
	if form.FormValid() != validChar {
		t.Error("Error Input Character must be 5 or greater than 5 @Valid_Len_Character")
	}
}

func TestForm_ValidPassword(t *testing.T) {
	rq := httptest.NewRequest("POST", "/", nil)
	form := NewForm(rq.PostForm)
	validPassword := form.ValidPassword("password", 15, rq)
	if validPassword {
		t.Error("Error input valid Password @Valid_password Testing")
	}
	//rq, _ = http.NewRequest("POST", "/", nil)
	postForm := url.Values{}
	postForm.Add("password", "123456789011111123")

	rq, _ = http.NewRequest("POST", "/", nil)
	//form = NewForm(rq.PostForm)

	rq.PostForm = postForm
	form = NewForm(rq.PostForm)

	validPassword = form.ValidPassword("password", 15, rq)
	if !validPassword && !form.FormValid() {
		t.Error("Error invalid Password @Valid_password Testing")
	}
}
