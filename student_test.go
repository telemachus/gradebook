package gradebook_test

import (
	"gradebook"
	"testing"
)

func TestNewStudentValid(t *testing.T) {
	t.Parallel()

	firstName := "First"
	lastName := "Last"
	_, err := gradebook.NewStudent(firstName, lastName)
	if err != nil {
		t.Errorf("gradebook.NewStudent(y, z) = %v; want no error", err)
	}
}
