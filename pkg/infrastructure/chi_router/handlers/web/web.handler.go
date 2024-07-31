package web

import (
	"chi_boilerplate/utils"
	"html/template"
	"net/http"
)

// HealthCheck returns status code 200.
func HealthCheck(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusOK)

	return nil
}

// GetAPIv1Doc returns the API v1 documentation.
func GetAPIv1Doc(w http.ResponseWriter, r *http.Request) error {
	tmpl, err := template.ParseFiles("./templates/doc_api_v1.gohtml")
	if err != nil {
		return utils.Err500(w, err, "Error when parsing HTML template", nil)
	}

	if err := tmpl.Execute(w, nil); err != nil {
		return utils.Err500(w, err, "Error when executing HTML template", nil)
	}

	return nil
}
