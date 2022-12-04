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

// Terms stores information about divisions of the year.
type Terms map[string]*Term

// Categories stores grading categories.
type Categories []string

// CategoriesPretty stores transformations from abbreviations to display names.
type CategoriesPretty map[string]string

// CategoryWeights stores maps of type => percentage for grades.
type CategoryWeights map[string]float64

// TypesToCategories stores maps of type => category for grades.
type TypesToCategories map[string]string

// Student stores information about students.
type Student struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	// I will need a map here to store grades, right?
}

// Students stores student information.
type Students map[string]*Student

type Class struct {
	Name              string `json:"name"`
	Terms             `json:"terms"`
	Categories        `json:"categories"`
	CategoriesPretty  `json:"categories_pretty"`
	CategoryWeights   `json:"category_weights"`
	TypesToCategories `json:"types_to_categories"`
	Students          `json:"students"`
}

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
