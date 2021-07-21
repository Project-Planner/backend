package web

import (
	"fmt"
	"html/template"
	"net/http"
)

var errTemplate *template.Template

func writeError(w http.ResponseWriter, msg string, code int) {
	w.WriteHeader(code)

	if errTemplate == nil {
		http.Error(w, msg, code)
		return
	}

	err := errTemplate.Execute(w, struct {
		DetailedError string
		StatusCode    string
		Status        string
	}{DetailedError: msg, StatusCode: fmt.Sprint(code), Status: http.StatusText(code)})
	if err != nil {
		fmt.Println(err)
	}
}
