package api

import (
	"deplagene/snippetbox"
	"html/template"
	"io/fs"
	"path/filepath"

	"github.com/gin-gonic/gin/render"
	"github.com/jackc/pgx/v5/pgtype"
)

type TemplateRenderer struct {
	Templates map[string]*template.Template
}

func (t *TemplateRenderer) Instance(name string, data any) render.Render {
	return render.HTML{
		Template: t.Templates[name],
		Data:     data,
	}
}

func humanDate(t pgtype.Timestamp) string {
	if !t.Valid {
		return ""
	}
	return t.Time.Format("02 Jan 2006 at 15:04")
}

func NewTemplateRenderer() (*TemplateRenderer, error) {
	templates := make(map[string]*template.Template)

	pages, err := fs.Glob(snippetbox.UIFS, "ui/html/*-page.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		funcs := template.FuncMap{
			"humanDate": humanDate,
		}

		ts, err := template.New(name).Funcs(funcs).ParseFS(snippetbox.UIFS, "ui/html/base-layout.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFS(snippetbox.UIFS, "ui/html/footer-partial.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFS(snippetbox.UIFS, page)
		if err != nil {
			return nil, err
		}

		templates[name] = ts
	}

	return &TemplateRenderer{Templates: templates}, nil
}
