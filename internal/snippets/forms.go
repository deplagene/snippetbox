package snippets

import (
	"strings"
	"unicode/utf8"
)

type SnippetCreateForm struct {
	Title   string            `form:"title"`
	Content string            `form:"content"`
	Expires string            `form:"expires"`
	Errors  map[string]string `form:"-"`
}

func (f *SnippetCreateForm) IsValid() bool {
	f.Errors = make(map[string]string)

	if strings.TrimSpace(f.Title) == "" {
		f.Errors["Title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(f.Title) > 100 {
		f.Errors["Title"] = "This field cannot be more than 100 characters long"
	}

	if strings.TrimSpace(f.Content) == "" {
		f.Errors["Content"] = "This field cannot be blank"
	}

	if strings.TrimSpace(f.Expires) == "" {
		f.Errors["Expires"] = "This field cannot be blank"
	} else if f.Expires != "365" && f.Expires != "7" && f.Expires != "1" {
		f.Errors["Expires"] = "This field must be 365, 7 or 1"
	}

	return len(f.Errors) == 0
}