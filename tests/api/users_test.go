package api

import (
	"chi_boilerplate/pkg/domain/requests"
	"chi_boilerplate/tests/helpers"
	"strings"
	"testing"
)

func TestUserLogin(t *testing.T) {
	tdb := helpers.InitMySQL("../../.env", "../../migrations")
	defer tdb.Drop()

	useCases := []helpers.Test{
		{
			Description: "User login",
			Route:       "/api/v1/token",
			Method:      "POST",
			Body: strings.NewReader(helpers.JsonToString(requests.UserLogin{
				Email:    helpers.UserEmail,
				Password: helpers.UserPassword,
			})),
			Headers: []helpers.Header{
				{Key: "Content-Type", Value: "application/json; charset=utf-8"},
			},
			CheckCode:    true,
			ExpectedCode: 200,
		},
	}

	tdb.Execute(t, useCases, "../../templates")
}

func TestUserCreation(t *testing.T) {
	tdb := helpers.InitMySQL("../../.env", "../../migrations")
	defer tdb.Drop()

	useCases := []helpers.Test{
		{
			Description: "User creation",
			Route:       "/api/v1/users",
			Method:      "POST",
			Body: strings.NewReader(helpers.JsonToString(requests.UserCreation{
				Email:     "test1@gmail.com",
				Password:  "11111111",
				Lastname:  "Test",
				Firstname: "Creation",
			})),
			Headers: []helpers.Header{
				{Key: "Content-Type", Value: "application/json; charset=utf-8"},
				{Key: "Authorization", Value: "Bearer " + tdb.Token},
			},
			CheckCode:    true,
			ExpectedCode: 200,
		},
		{
			Description: "User creation with invalid password",
			Route:       "/api/v1/users",
			Method:      "POST",
			Body: strings.NewReader(helpers.JsonToString(requests.UserCreation{
				Email:     "test1@gmail.com",
				Password:  "1111111",
				Lastname:  "Test",
				Firstname: "Creation",
			})),
			Headers: []helpers.Header{
				{Key: "Content-Type", Value: "application/json; charset=utf-8"},
				{Key: "Authorization", Value: "Bearer " + tdb.Token},
			},
			CheckCode:    true,
			CheckBody:    true,
			ExpectedCode: 400,
			ExpectedBody: `{"code":400,"message":"Invalid request data","details":[{"FailedField":"Password","Tag":"min","Value":"8"}]}`,
		},
		{
			Description: "User creation with invalid email",
			Route:       "/api/v1/users",
			Method:      "POST",
			Body: strings.NewReader(helpers.JsonToString(requests.UserCreation{
				Email:     "test1gmail.com",
				Password:  "11111111",
				Lastname:  "Test",
				Firstname: "Creation",
			})),
			Headers: []helpers.Header{
				{Key: "Content-Type", Value: "application/json; charset=utf-8"},
				{Key: "Authorization", Value: "Bearer " + tdb.Token},
			},
			CheckCode:    true,
			CheckBody:    true,
			ExpectedCode: 400,
			ExpectedBody: `{"code":400,"message":"Invalid request data","details":[{"FailedField":"Email","Tag":"email","Value":""}]}`,
		},
	}

	tdb.Execute(t, useCases, "../../templates")
}

func TestUserByID(t *testing.T) {
	tdb := helpers.InitMySQL("../../.env", "../../migrations")
	defer tdb.Drop()

	useCases := []helpers.Test{
		{
			Description: "Get user by ID",
			Route:       "/api/v1/users/" + helpers.UserID,
			Method:      "GET",
			Headers: []helpers.Header{
				{Key: "Content-Type", Value: "application/json; charset=utf-8"},
				{Key: "Authorization", Value: "Bearer " + tdb.Token},
			},
			CheckCode:    true,
			ExpectedCode: 200,
		},
		{
			Description: "Get user by ID with unknown user ID",
			Route:       "/api/v1/users/f47ac10b-58cc-0372-8562-0b8e853961b3",
			Method:      "GET",
			Headers: []helpers.Header{
				{Key: "Content-Type", Value: "application/json; charset=utf-8"},
				{Key: "Authorization", Value: "Bearer " + tdb.Token},
			},
			CheckCode:    true,
			ExpectedCode: 404,
		},
	}

	tdb.Execute(t, useCases, "../../templates")
}

func TestUserDelete(t *testing.T) {
	tdb := helpers.InitMySQL("../../.env", "../../migrations")
	defer tdb.Drop()

	useCases := []helpers.Test{
		{
			Description: "Delete user by ID",
			Route:       "/api/v1/users/" + helpers.UserID,
			Method:      "DELETE",
			Headers: []helpers.Header{
				{Key: "Content-Type", Value: "application/json; charset=utf-8"},
				{Key: "Authorization", Value: "Bearer " + tdb.Token},
			},
			CheckCode:    true,
			ExpectedCode: 204,
		},
		{
			Description: "Delete user with unknown user ID",
			Route:       "/api/v1/users/f47ac10b-58cc-0372-8562-0b8e853961b3",
			Method:      "DELETE",
			Headers: []helpers.Header{
				{Key: "Content-Type", Value: "application/json; charset=utf-8"},
				{Key: "Authorization", Value: "Bearer " + tdb.Token},
			},
			CheckCode:    true,
			ExpectedCode: 404,
		},
	}

	tdb.Execute(t, useCases, "../../templates")
}

func TestUserGetAll(t *testing.T) {
	tdb := helpers.InitMySQL("../../.env", "../../migrations")
	defer tdb.Drop()

	useCases := []helpers.Test{
		{
			Description: "Get all users without query parameters",
			Route:       "/api/v1/users/",
			Method:      "GET",
			Headers: []helpers.Header{
				{Key: "Content-Type", Value: "application/json; charset=utf-8"},
				{Key: "Authorization", Value: "Bearer " + tdb.Token},
			},
			CheckCode:    true,
			ExpectedCode: 200,
			CheckBody:    true,
			ExpectedBody: `{"data":[{"id":"` + helpers.UserID + `","email":"` + helpers.UserEmail + `","lastname":"Test","firstname":"Test","created_at":"` + helpers.UserCreatedAt + `","updated_at":"` + helpers.UserUpdatedAt + `"}],"total":1}`,
		},
	}

	tdb.Execute(t, useCases, "../../templates")
}
