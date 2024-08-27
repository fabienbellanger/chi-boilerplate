package values_objects

import (
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
func NewID() ID {
	return ID{Value: uuid.New()}
}

// NewID creates a new ID from string
func NewIDFrom(value string) (ID, error) {
	uid, err1 := uuid.Parse(value)
	if err1 != nil {
		return ID{}, err1
	}
	id := ID{Value: uid}

	return id, nil
}
