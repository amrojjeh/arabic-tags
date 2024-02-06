package validator

import (
	"fmt"
	"net/mail"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	field      string
	value      string
	errMessage string
	valid      bool
}

func NewValidator(field, value string) Validator {
	return Validator{
		field: field,
		value: value,
		valid: true,
	}
}

func (v Validator) SetError(message string) Validator {
	v.valid = false
	v.errMessage = message
	return v
}

func (v Validator) MaxLength(max int) Validator {
	if utf8.RuneCountInString(v.value) > max {
		v = v.SetError(fmt.Sprintf("Your %s exceeds the max length of %v",
			v.field, max))
	}
	return v
}

func (v Validator) MaxBytes(max int) Validator {
	if len(v.value) > max {
		v = v.SetError(fmt.Sprintf("Your %s exceeds the max length of %v",
			v.field, max))
	}
	return v
}

func (v Validator) Required() Validator {
	if v.value == "" {
		v = v.SetError(fmt.Sprintf("Please provide a %s", v.field))
	}
	return v
}

func (v Validator) SameAs(val string) Validator {
	if v.value != val {
		v = v.SetError(fmt.Sprintf("Your %s did not match", v.field))
	}
	return v
}

func (v Validator) IsEmail() Validator {
	_, err := mail.ParseAddress(v.value)
	if err != nil {
		v = v.SetError(fmt.Sprintf("Your %s is not a valid email address", v.field))
	}
	return v
}

func (v Validator) NotBlank() Validator {
	if strings.TrimSpace(v.value) == "" {
		v = v.SetError(fmt.Sprintf("You must provide a %s", v.field))
	}
	return v
}

func (v Validator) CustomMessage(message string) Validator {
	v.errMessage = message
	return v
}

func (v Validator) Validate() string {
	if !v.valid {
		return v.errMessage
	}
	return ""
}
