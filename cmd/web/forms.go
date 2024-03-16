package main

import (
	"unicode/utf8"

	"example.com/pkg/validator"
)

type RegistrationForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	ConfirmPassword     string `form:"confirm-password"`
	validator.Validator `form:"-"`
}

func (form *RegistrationForm) validate() {
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.IsEmail(form.Email), "email", "Please enter a valid email")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.Between(utf8.RuneCountInString(form.Password), 8, 16), "password", "This field should be between 8 & 16 characters")
	form.CheckField(validator.In(form.ConfirmPassword, form.Password), "confirm_password", "Password Mismatch")
}

func (form *RegistrationForm) isValid() bool {
	form.validate()

	return !form.Validator.HasErrors()
}

type LoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	Remember            bool   `form:"remember"`
	validator.Validator `form:"-"`
}

func (form *LoginForm) validate() {
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.IsEmail(form.Email), "email", "Please enter a valid email")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
}

func (form *LoginForm) isValid() bool {
	form.validate()

	return !form.Validator.HasErrors()
}
