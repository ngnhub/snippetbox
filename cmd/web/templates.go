package main

import (
	"net/url"
	"path/filepath"
	"text/template"
	"time"

	"github.com/ngnhub/snippetbox/pkg/models"
)

// Create a humanDate function which returns a nicely formatted string
// representation of a time.Time object.
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

type templateData struct {
	CurrentYear int
	FormData    url.Values
	FormErrors  map[string]string
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// essentially gives us a slice of all the 'page' templates
	path := filepath.Join(dir, "*.page.tmpl")
	pagePaths, err := filepath.Glob(path)
	if err != nil {
		return nil, err
	}

	for _, page := range pagePaths {
		// Extract the file name (like 'home.page.tmpl')
		pageName := filepath.Base(page)

		template, err := template.New(pageName).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}
		template, err = template.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}
		template, err = template.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}
		cache[pageName] = template
	}
	return cache, nil
}
