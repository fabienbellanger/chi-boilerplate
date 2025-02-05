package tests

import (
	"chi_boilerplate/tests/helpers"
	"strings"
	"testing"
)

func TestWebRoutes(t *testing.T) {
	tdb := helpers.InitMySQL("../.env", "../migrations")
	defer tdb.Drop()

	useCases := []helpers.Test{
		{
			Description:  "Health Check route",
			Route:        "/health",
			Method:       "GET",
			CheckCode:    true,
			CheckBody:    true,
			ExpectedCode: 200,
			ExpectedBody: "",
		},
		{
			Description: "Non existing route",
			Route:       "/not-exists",
			Method:      "GET",
			Body:        strings.NewReader("v=1"),
			Headers: []helpers.Header{
				{Key: "Content-Type", Value: "application/x-www-form-urlencoded"},
			},
			CheckCode:     true,
			CheckBody:     true,
			ExpectedError: false,
			ExpectedCode:  404,
			ExpectedBody:  "{\"code\":404,\"message\":\"Ressource not found\"}",
		},
	}

	tdb.Execute(t, useCases, "../templates")
}
