// Package gradebook is a library to read and write gradebook files.
package gradebook

import (
	"encoding/json"
	"os"
)

// Term stores dates for the start and end of a term.
type Term struct {
	Start string
	End   string
}

// Terms maps strings (e.g., "q1") to Term structs.
type Terms map[string]*Term

// Categories stores grading categories.
type Categories []string

// CategoriesPretty maps category names to a more readable representation.
type CategoriesPretty map[string]string

// CategoryWeights maps categories to their value in grading.
type CategoryWeights map[string]int

// TypesToCategories maps assignment types to categories for grading.
type TypesToCategories map[string]string

// Student stores information about students.
type Student struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	// I will need a map here to store grades, right?
}

// Students maps emails to Student structs.
//
// Email is an appropriate equivalent to a database's primary key because
// emails are unique.
type Students map[string]*Student

// Class stores information about the structure of a class and its students.
type Class struct {
	Name              string `json:"name"`
	Terms             `json:"terms"`
	Categories        `json:"categories"`
	CategoriesPretty  `json:"categories_pretty"`
	CategoryWeights   `json:"category_weights"`
	TypesToCategories `json:"types_to_categories"`
	Students          `json:"students"`
}

// LoadClass unmarshals a class.json file into a pointer to Class.
func LoadClass(classFile string) (*Class, error) {
	data, err := os.ReadFile(classFile)
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
