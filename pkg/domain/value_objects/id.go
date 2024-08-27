package values_objects

import (
	"chi_boilerplate/utils"

	"github.com/google/uuid"
)

type ID struct {
	Value uuid.UUID `validate:"required,uuid"`
}

// String returns the ID value
func (id *ID) String() string {
	return id.Value.String()
}

// NewID creates a new ID
func NewID() (ID, error) {
	id := ID{Value: uuid.New()}

	err := id.Validate()
	if err != nil {
		return ID{}, &err
	}

	return id, nil
}

// NewID creates a new ID from string
func NewIDFrom(value string) (ID, error) {
	id := ID{Value: uuid.MustParse(value)}

	err := id.Validate()
	if err != nil {
		return ID{}, &err
	}

	return id, nil
}

// Validate checks if a struct is valid and returns an array of errors
func (id *ID) Validate() utils.ValidatorErrors {
	return utils.ValidateStruct(id)
}
