package responses

import (
	"chi_boilerplate/pkg/domain/entities"
	values_objects "chi_boilerplate/pkg/domain/value_objects"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUserByIdToUserHTTP(t *testing.T) {
	id, _ := values_objects.NewIDFrom("f47ac10b-58cc-0372-8562-0b8e853961a1")
	email, _ := values_objects.NewEmail("test@test.com")
	tt, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
	user := UserById{
		ID:        id,
		Email:     email,
		Lastname:  "Doe",
		Firstname: "John",
		CreatedAt: tt,
		UpdatedAt: tt,
	}

	expected := UserHTTP{
		ID:        "f47ac10b-58cc-0372-8562-0b8e853961a1",
		Email:     "test@test.com",
		Lastname:  "Doe",
		Firstname: "John",
		CreatedAt: "2021-01-01T00:00:00Z",
		UpdatedAt: "2021-01-01T00:00:00Z",
	}

	assert.Equal(t, user.ToUserHTTP(), expected)
}

func TestUserCreationToUserHTTP(t *testing.T) {
	id, _ := values_objects.NewIDFrom("f47ac10b-58cc-0372-8562-0b8e853961a1")
	email, _ := values_objects.NewEmail("test@test.com")
	tt, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
	user := UserCreation{
		ID:        id,
		Email:     email,
		Lastname:  "Doe",
		Firstname: "John",
		CreatedAt: tt,
		UpdatedAt: tt,
	}

	expected := UserHTTP{
		ID:        "f47ac10b-58cc-0372-8562-0b8e853961a1",
		Email:     "test@test.com",
		Lastname:  "Doe",
		Firstname: "John",
		CreatedAt: "2021-01-01T00:00:00Z",
		UpdatedAt: "2021-01-01T00:00:00Z",
	}

	assert.Equal(t, user.ToUserHTTP(), expected)
}

func TestUserLoginRepository(t *testing.T) {
	id, _ := values_objects.NewIDFrom("f47ac10b-58cc-0372-8562-0b8e853961a1")
	email, _ := values_objects.NewEmail("test@test.com")
	tt, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
	user := UserLoginRepository{
		ID:        id.String(),
		Email:     "test@test.com",
		Lastname:  "Doe",
		Firstname: "john",
		CreatedAt: "2021-01-01T00:00:00Z",
		UpdatedAt: "2021-01-01T00:00:00Z",
	}

	expected := entities.User{
		ID:        id,
		Email:     email,
		Lastname:  "Doe",
		Firstname: "john",
		CreatedAt: tt,
		UpdatedAt: tt,
	}
	got, _ := user.ToUser()

	assert.Equal(t, got, expected)
}
