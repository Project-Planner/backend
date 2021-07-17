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
		detailedError string
		code string
		status string
	}{detailedError: msg, code: fmt.Sprint(code), status: http.StatusText(code)})
}
