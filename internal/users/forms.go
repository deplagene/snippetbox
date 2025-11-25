package users

import (
	"regexp"
	"strings"
)

var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

type RegisterForm struct {
	Name     string
	Email    string
	Password string
	Errors   map[string]string
}

func (f *RegisterForm) IsValid() bool {
	f.Errors = make(map[string]string)

	if strings.TrimSpace(f.Name) == "" {
		f.Errors["Name"] = "This field cannot be blank"
	}

	if strings.TrimSpace(f.Email) == "" {
		f.Errors["Email"] = "This field cannot be blank"
	} else if !EmailRX.MatchString(f.Email) {
		f.Errors["Email"] = "This field must be a valid email address"
	}

	if strings.TrimSpace(f.Password) == "" {
		f.Errors["Password"] = "This field cannot be blank"
	} else if len(f.Password) < 8 {
		f.Errors["Password"] = "This field must be at least 8 characters long"
	}

	return len(f.Errors) == 0
}

type LoginForm struct {
	Email    string
	Password string
	Errors   map[string]string
}

func (f *LoginForm) IsValid() bool {
	f.Errors = make(map[string]string)

	if strings.TrimSpace(f.Email) == "" {
		f.Errors["Email"] = "This field cannot be blank"
	} else if !EmailRX.MatchString(f.Email) {
		f.Errors["Email"] = "This field must be a valid email address"
	}

	if strings.TrimSpace(f.Password) == "" {
		f.Errors["Password"] = "This field cannot be blank"
	}

	return len(f.Errors) == 0
}