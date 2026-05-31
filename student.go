package gradebook

import "errors"

// NewStudent a new *Student. If firstName or lastName is empty, the method
// returns an error.
func NewStudent(firstName, lastName string) (*Student, error) {
	if firstName == "" {
		return nil, errors.New("gradebook: first name cannot be empty")
	}
	if lastName == "" {
		return nil, errors.New("gradebook: last name cannot be empty")
	}

	return &Student{firstName: firstName, lastName: lastName}, nil
}
