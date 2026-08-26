package gradebook

import (
	"encoding/json/v2"
	"fmt"
	"os"
	"path/filepath"
)

type classInitMode uint8

const (
	initGrades classInitMode = 1 << iota
	initUnscored
)

// ParseClassFileForUnscored parses a class.json file into a pointer to Class.
// Unlike ParseClassFile, this function creates the UnscoredByCategory map
// needed to count unscored assignments.
func ParseClassFileForUnscored(classFile string) (*Class, error) {
	return parseClassFileWithInit(classFile, initUnscored)
}

// ParseClassFileForGrades parses a class.json file into a pointer to Class.
// Unlike ParseClassFile, this function creates the GradesByCategory map needed
// to store grades.
func ParseClassFileForGrades(classFile string) (*Class, error) {
	return parseClassFileWithInit(classFile, initGrades)
}

func parseClassFileWithInit(classFile string, mode classInitMode) (*Class, error) {
	data, err := os.ReadFile(filepath.Clean(classFile))
	if err != nil {
		return nil, fmt.Errorf("gradebook: read class file %q: %w", classFile, err)
	}

	var spec ClassSpec
	err = json.Unmarshal(data, &spec, json.DefaultOptionsV2())
	if err != nil {
		return nil, fmt.Errorf("gradebook: unmarshal class file %q: %w", classFile, err)
	}

	class := cloneClassFromSpec(&spec)
	if err := class.trustParsedData(); err != nil {
		return nil, &InvalidClassError{Err: err}
	}

	class.initializeStudentMaps(mode)

	return class, nil
}

func (c *Class) initializeStudentMaps(mode classInitMode) {
	for _, student := range c.studentsForInitialization() {
		if student == nil {
			continue
		}

		c.initializeStudentGradesMap(student, mode)
		c.initializeStudentUnscoredMap(student, mode)
	}

	c.projectTrustedStudentsToCompatibilityFields()
}

func (c *Class) studentsForInitialization() map[string]*Student {
	if domain := c.trustedDomain(); domain != nil {
		return domain.studentsByEmail
	}

	return c.studentsByEmail
}

func (c *Class) assignmentCategoriesForInitialization() []string {
	if domain := c.trustedDomain(); domain != nil {
		return domain.assignmentCategories
	}

	return c.assignmentCategories
}

func (c *Class) initializeStudentGradesMap(student *Student, mode classInitMode) {
	if mode&initGrades == 0 {
		return
	}

	categories := c.assignmentCategoriesForInitialization()
	if student.gradesByCategory == nil {
		student.gradesByCategory = make(map[string][]float64, len(categories))
	}
	for _, cat := range categories {
		if _, ok := student.gradesByCategory[cat]; ok {
			continue
		}
		student.gradesByCategory[cat] = make([]float64, 0, 25)
	}
}

func (c *Class) initializeStudentUnscoredMap(student *Student, mode classInitMode) {
	if mode&initUnscored == 0 {
		return
	}

	categories := c.assignmentCategoriesForInitialization()
	if student.unscoredByCategory == nil {
		student.unscoredByCategory = make(map[string]int, len(categories))
	}
	for _, cat := range categories {
		if _, ok := student.unscoredByCategory[cat]; ok {
			continue
		}
		student.unscoredByCategory[cat] = 0
	}
}
