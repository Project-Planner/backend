//go:generate go run routes_generator.go

package main

import (
	"os"
	"strings"
	"text/template"
)

const tmpl = `
// AUTO-GENERATED CODE; DO NOT EDIT

package web

import (
	"fmt"
	"github.com/gorilla/mux"
)

func attachEndpoints(r *mux.Router) {
{{- $methods := .Methods -}}
{{- $items := .Items -}}
{{- range $idxI, $item := $items}}
	{{lowerCasePlural $item}}Router := r.PathPrefix("/api/{{lowerCasePlural $item}}").Subrouter()

	{{lowerCasePlural $item}}Router.HandleFunc(fmt.Sprintf("/post/{%s}/{%s}", userIDStr, calendarIDStr), post{{$item}}Handler).Methods("POST")
	{{lowerCasePlural $item}}Router.HandleFunc(fmt.Sprintf("/post/{%s}", calendarIDStr), post{{$item}}Handler).Methods("POST")
	{{lowerCasePlural $item}}Router.HandleFunc("/post", post{{$item}}Handler).Methods("POST")

	{{lowerCasePlural $item}}Router.HandleFunc(fmt.Sprintf("/other/{%s}/{%s}/{%s}", userIDStr, calendarIDStr, itemIDStr), methodHandler(nil
{{- range $idxM, $method := $methods -}}
{{- if isPost $method -}}
{{- else -}}
	, {{lowerCase $method}}{{$item}}Handler
{{- end -}}
{{- end -}}
	)).Methods("POST")	
	{{lowerCasePlural $item}}Router.HandleFunc(fmt.Sprintf("/other/{%s}/{%s}", calendarIDStr, itemIDStr), methodHandler(nil 
{{- range $idxM, $method := $methods -}}
{{- if isPost $method -}}
{{- else -}}
	, {{lowerCase $method}}{{$item}}Handler
{{- end -}}
{{- end -}}
	)).Methods("POST")
	{{lowerCasePlural $item}}Router.HandleFunc(fmt.Sprintf("/other/{%s}", itemIDStr), methodHandler(nil 
{{- range $idxM, $method := $methods -}}
{{- if isPost $method -}}
{{- else -}}
	, {{lowerCase $method}}{{$item}}Handler
{{- end -}}
{{- end -}}
	)).Methods("POST")

{{ end }}
}
`

func main() {
	f, err := os.Create("../endpoints.go")
	if err != nil {
		panic(err)
	}

	fm := template.FuncMap{
		"lowerCasePlural": lowerCasePlural,
		"isPost":          isPost,
		"lowerCase":       lowerCase,
	}

	err = template.Must(template.New("").Funcs(fm).Parse(tmpl)).Execute(f, struct {
		Items   []string
		Methods []string
	}{
		Items: []string{
			"Appointment",
			"Milestone",
			"Task",
		},
		Methods: []string{
			"POST",
			"PUT",
			"DELETE",
		},
	})
	if err != nil {
		panic(err)
	}
}

func lowerCasePlural(s string) string {
	return strings.ToLower(s) + "s"
}

func lowerCase(s string) string {
	return strings.ToLower(s)
}

func isPost(s string) bool {
	return s == "POST"
}
