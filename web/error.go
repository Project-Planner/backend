package web

import (
	"fmt"
	"html/template"
	"net/http"
)

var errTemplate *template.Template

func writeError(w http.ResponseWriter, msg string, code int) {
	w.WriteHeader(code)

	_ = errTemplate.Execute(w, struct {
		DetailedError string
		Code          string
		Status        string
	}{DetailedError: msg, Code: fmt.Sprint(code), Status: http.StatusText(code)})
}
