// Package gradebook is a library to read and write gradebook files.
package gradebook

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
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

// AssignmentCategories stores basic assignment categories, the broad groupings
// such as "major" and "minor".
type AssignmentCategories []string

// LabelsByAssignmentCategory maps human-readable labels by a category of
// assignment. E.g., the category "cp" has the label "Class Participation",
// and the category "major" has the label "Major Assessments".
type LabelsByAssignmentCategory map[string]string

// WeightsByAssignmentCategory maps percentage values in a grading rubric by
// assignment category. The sum of the weights must equal 100 in order for this
// type to be valid.
type WeightsByAssignmentCategory map[string]int

// CategoriesByAssignmentType maps categories by their assignment type. (E.g.,
// "test", "essay", and "project" all have the category "major". Every
// assignment type must belong to one and only one category, and every
// assignment type must be present in CategoriesByAssignmentType. Also, and
// this is less obvious, every assignment category must have an assignment
// type. If a category has only a single assignment type, the category and type
// will often have the same name.  E.g., both the category and the type are
// "cp".
type CategoriesByAssignmentType map[string]string

// AssignmentRecord represents a grade for a particular student on a particular
// assignment. If AssignmentRecord.Grade is nil, then no score has been
// recorded yet and the record should not be counted when calculating grades.
//
//nolint:govet // JSON field order matters more here than memory alignment.
type AssignmentRecord struct {
	Email string   `json:"email"`
	Grade *float64 `json:"grade"`
}

// AssignmentRecords stores AssignmentRecord structs.
type AssignmentRecords []*AssignmentRecord

// Gradebook represents a single gradebook file.
type Gradebook struct {
	AssignmentDate     string            `json:"assignment_date"`
	AssignmentName     string            `json:"assignment_name"`
	AssignmentType     string            `json:"assignment_type"`
	AssignmentCategory string            `json:"assignment_category"`
	AssignmentRecords  AssignmentRecords `json:"assignment_records"`
}

// Student represents a student.
type Student struct {
	GradesByCategory   map[string][]float64
	UnscoredByCategory map[string]int
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
}

// StudentsByEmail maps students by their email. (NB: an email is an
// appropriate equivalent to a database's primary key because emails are
// unique.)
type StudentsByEmail map[string]*Student

// Class represents a class and its students.
type Class struct {
	TermsByID                   `json:"terms_by_id"`
	LabelsByAssignmentCategory  `json:"labels_by_assignment_category"`
	WeightsByAssignmentCategory `json:"weights_by_assignment_category"`
	CategoriesByAssignmentType  `json:"categories_by_assignment_type"`
	StudentsByEmail             `json:"students_by_email"`
	Name                        string `json:"name"`
	AssignmentCategories        `json:"assignment_categories"`
}

// UnmarshalClass unmarshals a class.json file into a pointer to Class.
func UnmarshalClass(classFile string) (*Class, error) {
	return unmarshalClassWithInit(classFile, 0)
}

// UnmarshalGradebook unmarshals a gradebook file into a pointer to Gradebook.
func UnmarshalGradebook(gradebookFile string) (*Gradebook, error) {
	data, err := os.ReadFile(filepath.Clean(gradebookFile))
	if err != nil {
		return nil, fmt.Errorf("gradebook: read gradebook file %q: %w", gradebookFile, err)
	}

	var gradebook Gradebook
	err = json.Unmarshal(data, &gradebook)
	if err != nil {
		return nil, fmt.Errorf("gradebook: unmarshal gradebook file %q: %w", gradebookFile, err)
	}

	return &gradebook, nil
}

// dateSnip gets the date string from the end of a gradebook filename. If the
// function does not find a valid date (in YYYYMMDD format), then it return an
// error.
func dateSnip(name string) (string, error) {
	base := filepath.Base(name)
	stem := strings.TrimSuffix(base, gradebookSuffix)
	if len(stem) < dateFmtLen {
		return "", fmt.Errorf("gradebook: invalid yyyymmdd date in gradebook file name %q", base)
	}

	dateStr := stem[len(stem)-dateFmtLen:]
	if _, err := time.Parse("20060102", dateStr); err != nil {
		return "", fmt.Errorf("gradebook: invalid yyyymmdd date in gradebook file name %q", base)
	}

	return dateStr, nil
}

// LoadGrades scans a given directory for *.gradebook files and adds grades
// from those files to students. The method returns an error if there is
// a problem reading, unmarshaling, or closing a file.
func (c *Class) LoadGrades(dir string, term *Term) error {
	c.initializeStudentMaps(initGrades)

	return c.loadGradebooks(dir, term, false)
}

// LoadUnscored scans a given directory for *.gradebook files and counts
// unscored assignments for each student by assignment category. The method
// returns an error if there is a problem reading, unmarshaling, or closing a
// file.
func (c *Class) LoadUnscored(dir string, term *Term) error {
	c.initializeStudentMaps(initUnscored)

	return c.loadGradebooks(dir, term, true)
}

func (c *Class) loadGradebooks(dir string, term *Term, countUnscored bool) error {
	gradebooks, err := gradebookFilesInDir(dir)
	if err != nil {
		return err
	}

	for _, gradebook := range gradebooks {
		if term != nil {
			dateStr, err := dateSnip(gradebook)
			if err != nil {
				return err
			}

			if !term.Includes(dateStr) {
				continue
			}
		}

		if err := c.loadGradebookFile(gradebook, countUnscored); err != nil {
			return err
		}
	}

	return nil
}

func gradebookFilesInDir(dir string) ([]string, error) {
	entries, err := os.ReadDir(filepath.Clean(dir))
	if err != nil {
		return nil, fmt.Errorf("gradebook: read directory %q: %w", dir, err)
	}

	gradebooks := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if filepath.Ext(entry.Name()) != gradebookSuffix {
			continue
		}

		gradebooks = append(gradebooks, filepath.Join(dir, entry.Name()))
	}

	return gradebooks, nil
}

func (c *Class) loadGradebookFile(gradebookPath string, countUnscored bool) error {
	gbData, err := UnmarshalGradebook(gradebookPath)
	if err != nil {
		return err
	}

	category, err := c.categoryForAssignmentType(gbData.AssignmentType)
	if err != nil {
		return err
	}

	for i, ar := range gbData.AssignmentRecords {
		if ar == nil {
			return fmt.Errorf("gradebook: nil assignment record at index %d in %q", i, gradebookPath)
		}

		student, err := c.studentByEmail(ar.Email)
		if err != nil {
			return err
		}

		if countUnscored {
			if ar.Grade != nil {
				continue
			}
			student.UnscoredByCategory[category]++

			continue
		}

		if ar.Grade == nil {
			continue
		}

		_, ok := student.GradesByCategory[category]
		if !ok {
			return fmt.Errorf("gradebook: unrecognized assignment category %q for type %q", category, gbData.AssignmentType)
		}

		student.GradesByCategory[category] = append(student.GradesByCategory[category], *ar.Grade)
	}

	return nil
}

func (c *Class) studentByEmail(email string) (*Student, error) {
	student, ok := c.StudentsByEmail[email]
	if !ok {
		return nil, fmt.Errorf("gradebook: no student with email %q", email)
	}
	if student == nil {
		return nil, fmt.Errorf("gradebook: student with email %q is nil", email)
	}

	return student, nil
}

func (c *Class) categoryForAssignmentType(assignmentType string) (string, error) {
	category, ok := c.CategoriesByAssignmentType[assignmentType]
	if !ok {
		return "", fmt.Errorf("gradebook: unrecognized assignment type %q", assignmentType)
	}

	return category, nil
}
