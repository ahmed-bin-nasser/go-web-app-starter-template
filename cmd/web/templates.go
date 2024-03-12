package main

import (
	"html/template"
	"strings"

	"example.com/assets"
	"example.com/pkg/funcs"
)

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := assets.EmbeddedFiles.ReadDir("templates/pages")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := strings.Split(page.Name(), ".")[0]
		files := []string{"templates/base.go.tmpl", "templates/partials/*.go.tmpl", "templates/pages/" + page.Name()}

		ts, err := template.New("").Funcs(funcs.TemplateFuncs).ParseFS(assets.EmbeddedFiles, files...)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}

	return cache, nil
}
