package gradebook

import (
	"testing"
)

func TestNewStudentValid(t *testing.T) {
	t.Parallel()

	firstName := "First"
	lastName := "Last"
	_, err := NewStudent(firstName, lastName)
	if err != nil {
		t.Errorf("NewStudent(y, z) = %v; want no error", err)
	}
}
