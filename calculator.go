package gradebook

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type classInitMode uint8

const (
	initGrades classInitMode = 1 << iota
	initUnscored
)

// UnmarshalUnscoredClass unmarshals a class.json file into a pointer to Class.
// Unlike UnmarshalClass, this function creates the UnscoredByCategory map
// needed to count unscored assignments.
func UnmarshalUnscoredClass(classFile string) (*Class, error) {
	return unmarshalClassWithInit(classFile, initUnscored)
}

// UnmarshalCalcClass unmarshals a class.json file into a pointer to Class.
// Unlike UnmarshalClass, this function creates the Grades map needed to store
// grades.
func UnmarshalCalcClass(classFile string) (*Class, error) {
	return unmarshalClassWithInit(classFile, initGrades)
}

func unmarshalClassWithInit(classFile string, mode classInitMode) (*Class, error) {
	data, err := os.ReadFile(filepath.Clean(classFile))
	if err != nil {
		return nil, fmt.Errorf("gradebook: read class file %q: %w", classFile, err)
	}

	var class Class
	err = json.Unmarshal(data, &class)
	if err != nil {
		return nil, fmt.Errorf("gradebook: unmarshal class file %q: %w", classFile, err)
	}

	class.initializeStudentMaps(mode)

	return &class, nil
}

func (c *Class) initializeStudentMaps(mode classInitMode) {
	for _, student := range c.StudentsByEmail {
		if student == nil {
			continue
		}

		c.initializeStudentGradesMap(student, mode)
		c.initializeStudentUnscoredMap(student, mode)
	}
}

func (c *Class) initializeStudentGradesMap(student *Student, mode classInitMode) {
	if mode&initGrades == 0 {
		return
	}

	if student.GradesByCategory == nil {
		student.GradesByCategory = make(map[string][]float64, len(c.AssignmentCategories))
	}
	for _, cat := range c.AssignmentCategories {
		if _, ok := student.GradesByCategory[cat]; ok {
			continue
		}
		student.GradesByCategory[cat] = make([]float64, 0, 25)
	}
}

func (c *Class) initializeStudentUnscoredMap(student *Student, mode classInitMode) {
	if mode&initUnscored == 0 {
		return
	}

	if student.UnscoredByCategory == nil {
		student.UnscoredByCategory = make(map[string]int, len(c.AssignmentCategories))
	}
	for _, cat := range c.AssignmentCategories {
		if _, ok := student.UnscoredByCategory[cat]; ok {
			continue
		}
		student.UnscoredByCategory[cat] = 0
	}
}
