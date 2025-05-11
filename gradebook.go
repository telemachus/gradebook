// Package gradebook is a library to read and write gradebook files.
package gradebook

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	gradebookSuffix = ".gradebook"
	dateFmtLen      = len("YYYYMMDD")
)

// Term represents a grading period (e.g., a quarter or semester).
type Term struct {
	Start string
	End   string
}

// TermsByID maps Terms by short ID (e.g., "q1" points to the first quarter).
type TermsByID map[string]*Term

// AssignmentTypes stores basic assignment types.
type AssignmentTypes []string

// LabelsByAssignmentType maps human-readable labels by a type of assignment.
// E.g., the type "cp" has the label "Class Participation", and the type
// "major" has the label "Major Assessments".
type LabelsByAssignmentType map[string]string

// WeightsByAssignmentType maps percentage values in a grading rubric by
// assignment type. The sum of the weights must equal 100 in order for this
// type to be valid.
type WeightsByAssignmentType map[string]int

// CategoriesByAssignmentType maps assignment categories by their basic type.
// (E.g., "test", "essay", and "project" are all categories of the "major"
// assignment type.) Every category must belong to one and only one type, and
// every category must be present in CategoriesByAssignmentType. Also, and this
// is less obvious, every assignment type must have a category. If a type has
// only a single category, the category and type will often have the same name.
// E.g., "cp" is the category of the "cp" assignment type.
type CategoriesByAssignmentType map[string]string

// Grade represents a single grade
type Grade struct {
	Grade *float64
	Email string
}

// Grades stores Grade structs.
type Grades []*Grade

// Gradebook represents a single gradebook file.
type Gradebook struct {
	AssignmentDate     string `json:"assignment_date"`
	AssignmentType     string `json:"assignment_type"`
	AssignmentCategory string `json:"assignment_category"`
	Grades             `json:"assignment_grades"`
}

// Student represents a student.
type Student struct {
	GradesByType map[string][]*float64
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
}

// StudentsByEmail maps students by their email. (NB: an email is an
// appropriate equivalent to a database's primary key because emails are
// unique.)
type StudentsByEmail map[string]*Student

// Class represents a class and its students.
type Class struct {
	TermsByID                  `json:"terms_by_id"`
	LabelsByAssignmentType     `json:"labels_by_assignment_type"`
	WeightsByAssignmentType    `json:"weights_by_assignment_type"`
	CategoriesByAssignmentType `json:"categories_by_assignment_type"`
	StudentsByEmail            `json:"students_by_email"`
	Name                       string `json:"name"`
	AssignmentTypes            `json:"assignment_types"`
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

// UnmarshalGradebook unmarshals a gradebook file into a pointer to Gradebook.
func UnmarshalGradebook(gradebookFile string) (*Gradebook, error) {
	data, err := os.ReadFile(filepath.Clean(gradebookFile))
	if err != nil {
		return nil, err
	}

	var gradebook Gradebook
	err = json.Unmarshal(data, &gradebook)
	if err != nil {
		return nil, err
	}

	return &gradebook, nil
}

// dateSnip gets the date string from the end of a gradebook filename. If the
// function does not find a valid date (in YYYYMMDD format), then it return an
// error.
func dateSnip(dateStr string) (string, error) {
	dateStr = strings.TrimSuffix(dateStr, gradebookSuffix)
	dateStrLen := utf8.RuneCountInString(dateStr)
	if dateStrLen < dateFmtLen {
		return "", fmt.Errorf("[%s] does not contain a valid YYYYMMDD date", dateStr)
	}

	dateStr = dateStr[dateStrLen-dateFmtLen:]
	if _, err := time.Parse("20060102", dateStr); err != nil {
		return "", fmt.Errorf("[%s] does not contain a valid YYYYMMDD date", dateStr)
	}

	return dateStr, nil
}
