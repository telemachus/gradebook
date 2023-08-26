// Package gradebook is a library to read and write gradebook files.
package gradebook

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Term stores dates for the start and end of a term.
type Term struct {
	Start string
	End   string
}

// Terms associates Term structs with short labels (e.g., "q1" points to the
// first quarter.
type Terms map[string]*Term

// Categories stores grading categories.
type Categories []string

// PrettyCategories associates items in Categories with a form ready for
// display. E.g., the category "cp" points to "Class Participation", and
// "major" points to "Major Assessments".
type PrettyCategories map[string]string

// Weights associates items in Categories with their (percentage) value in
// a grading rubric. The sum of the weights must equal 100 in order for this
// type to be valid.
type Weights map[string]int

// Subcategories associates subcategories with items in Categories. (E.g.,
// "test", "essay", and "project" are all subcategories of the "major" grading
// category.) Every subcategory must belong to one and only one category, and
// every category must be present in Categories for this type to be valid.
// Also, and this is less obvious, every category must have a subcategory. If
// a category has only a single member, the subcategory and category will often
// have the same name. E.g., "cp" is a subcategory of "cp".
type Subcategories map[string]string

// Grade represents a student's single grade
type Grade struct {
	Email string
	Score *float64
}

// Grades stores a group of Grade structs.
type Grades []Grade

// Student stores information about students.
type Student struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	// I will need a map of some kind here to store grades.
}

// Students associates emails with Student structs. (NB: an email is an
// appropriate equivalent to a database's primary key because emails are
// unique.)
type Students map[string]*Student

// Class stores information about the structure of a class and its students.
type Class struct {
	Name             string `json:"name"`
	Terms            `json:"terms"`
	Categories       `json:"categories"`
	PrettyCategories `json:"pretty_categories"`
	Weights          `json:"weights"`
	Subcategories    `json:"subcategories"`
	Students         `json:"students"`
}

// Gradebook stores information about a single gradebook file.
type Gradebook struct {
	Grades `json:"assignment_grades"`
}

// UnmarshalClass unmarshals a class.json file into a pointer to Class.
func UnmarshalClass(classFile string) (*Class, error) {
	data, err := os.ReadFile(filepath.Clean(classFile))
	if err != nil {
		return nil, err
	}

	var class Class
	err = json.Unmarshal(data, &class)

	if err != nil {
		return nil, err
	}

	return &class, nil
}
