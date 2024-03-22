package response

import (
	"bytes"
	"html/template"
	"net/http"

	"example.com/assets"
	"example.com/pkg/funcs"
)

func Page(w http.ResponseWriter, status int, data any, pagePath string) error {
	return PageWithHeaders(w, status, data, nil, pagePath)
}

func PageWithHeaders(w http.ResponseWriter, status int, data any, headers http.Header, pagePath string) error {
	patterns := []string{"base.go.tmpl", "partials/*.go.tmpl", pagePath}

	return NamedTemplateWithHeaders(w, status, data, headers, "base", patterns...)
}

func NamedTemplate(w http.ResponseWriter, status int, data any, templateName string, patterns ...string) error {
	return NamedTemplateWithHeaders(w, status, data, nil, templateName, patterns...)
}

func NamedTemplateWithHeaders(w http.ResponseWriter, status int, data any, headers http.Header, templateName string, patterns ...string) error {
	for i := range patterns {
		patterns[i] = "templates/" + patterns[i]
	}

	ts, err := template.New("").Funcs(funcs.TemplateFuncs).ParseFS(assets.TemplateFiles, patterns...)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)

	err = ts.ExecuteTemplate(buf, templateName, data)
	if err != nil {
		return err
	}

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.WriteHeader(status)
	buf.WriteTo(w)

	return nil
}
