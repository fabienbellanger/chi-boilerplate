package tests

import (
	"strings"
	"testing"
)

func TestWebRoutes(t *testing.T) {
	tdb := Init("../.env")
	defer tdb.Drop()

	useCases := []Test{
		{
			Description:  "Health Check route",
			Route:        "/health-check",
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
			Headers: []Header{
				{Key: "Content-Type", Value: "application/x-www-form-urlencoded"},
			},
			CheckCode:     true,
			CheckBody:     true,
			ExpectedError: false,
			ExpectedCode:  404,
			ExpectedBody:  "{\"code\":404,\"message\":\"Ressource not found\"}",
		},
	}

	Execute(t, tdb.DB, useCases, "../templates")
}
