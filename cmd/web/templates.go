package main

import (
	"path/filepath"
	"text/template"

	"github.com/ngnhub/snippetbox/pkg/models"
)

type templateData struct {
	CurrentYear int
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

		template, err := template.ParseFiles(page)
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
