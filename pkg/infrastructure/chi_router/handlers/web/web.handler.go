package web

import (
	"chi_boilerplate/utils"
	"html/template"
	"net/http"
)

// HealthCheck returns status code 200.
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// GetAPIv1Doc returns the API v1 documentation.
func GetAPIv1Doc(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/doc_api_v1.gohtml")
	if err != nil {
		utils.Err500(w, err, "Error when parsing HTML template", nil)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		utils.Err500(w, err, "Error when executing HTML template", nil)
		return
	}
}
