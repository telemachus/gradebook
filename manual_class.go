package gradebook

import (
	"errors"
	"fmt"
)

// InitMode controls optional map initialization when building manual classes.
type InitMode uint8

// InitNone performs no student map initialization.
const (
	InitNone InitMode = 0
)

// InitGrades and InitUnscored control student map initialization.
const (
	InitGrades InitMode = 1 << iota
	InitUnscored
)

func validateInitMode(mode InitMode) error {
	const validModes = InitGrades | InitUnscored
	if mode&^validModes != 0 {
		return fmt.Errorf("gradebook: invalid init mode %d", mode)
	}

	return nil
}

// BuildClass validates and initializes a manual class spec.
func BuildClass(spec *ClassSpec, mode InitMode) (*Class, error) {
	if err := validateInitMode(mode); err != nil {
		return nil, err
	}
	if spec == nil {
		return nil, &InvalidClassError{Err: errors.New("gradebook: class is nil")}
	}

	class := cloneClassFromSpec(spec)
	if err := class.trustParsedData(); err != nil {
		return nil, &InvalidClassError{Err: err}
	}

	class.initializeStudentMaps(classInitMode(mode))

	return class, nil
}
