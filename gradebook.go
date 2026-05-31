// Package gradebook is a library to read and write gradebook files.
package gradebook

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
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
type AssignmentRecord struct {
	Grade *float64 `json:"grade"`
	Email string   `json:"email"`
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
	gradesByCategory   map[string][]float64
	unscoredByCategory map[string]int
	firstName          string
	lastName           string
}

// StudentsByEmail maps students by their email. (NB: an email is an
// appropriate equivalent to a database's primary key because emails are
// unique.)
type StudentsByEmail map[string]*Student

// Class represents a class and its students.
type Class struct {
	termsByID                   TermsByID
	labelsByAssignmentCategory  LabelsByAssignmentCategory
	weightsByAssignmentCategory WeightsByAssignmentCategory
	categoriesByAssignmentType  CategoriesByAssignmentType
	studentsByEmail             StudentsByEmail
	domain                      *trustedClassDomain
	name                        string
	assignmentCategories        AssignmentCategories
	trusted                     bool
}

type trustedClassDomain struct {
	termsByID            map[string]*Term
	categoriesByType     map[string]string
	labelsByCategory     map[string]string
	weightsByCategory    map[string]int
	studentsByEmail      map[string]*Student
	assignmentCategories []string
}

func cloneTermsByID(terms TermsByID) map[string]*Term {
	if terms == nil {
		return nil
	}

	cloned := make(map[string]*Term, len(terms))
	for id, term := range terms {
		if term == nil {
			cloned[id] = nil

			continue
		}

		cloned[id] = &Term{
			Start: term.Start,
			End:   term.End,
		}
	}

	return cloned
}

func cloneGradesByCategory(gradesByCategory map[string][]float64) map[string][]float64 {
	if gradesByCategory == nil {
		return nil
	}

	cloned := make(map[string][]float64, len(gradesByCategory))
	for category, grades := range gradesByCategory {
		cloned[category] = slices.Clone(grades)
	}

	return cloned
}

func cloneStudent(student *Student) *Student {
	if student == nil {
		return nil
	}

	return &Student{
		gradesByCategory:   cloneGradesByCategory(student.gradesByCategory),
		unscoredByCategory: maps.Clone(student.unscoredByCategory),
		firstName:          student.firstName,
		lastName:           student.lastName,
	}
}

func cloneStudentsByEmail(students map[string]*Student) map[string]*Student {
	if students == nil {
		return nil
	}

	cloned := make(map[string]*Student, len(students))
	for email, student := range students {
		cloned[email] = cloneStudent(student)
	}

	return cloned
}

func newTrustedClassDomain(c *Class) *trustedClassDomain {
	if c == nil {
		return nil
	}

	return &trustedClassDomain{
		termsByID:            cloneTermsByID(c.termsByID),
		assignmentCategories: slices.Clone(c.assignmentCategories),
		categoriesByType:     maps.Clone(c.categoriesByAssignmentType),
		labelsByCategory:     maps.Clone(c.labelsByAssignmentCategory),
		weightsByCategory:    maps.Clone(c.weightsByAssignmentCategory),
		studentsByEmail:      cloneStudentsByEmail(c.studentsByEmail),
	}
}

// ParseClassFile parses a class.json file into a pointer to Class.
func ParseClassFile(classFile string) (*Class, error) {
	return parseClassFileWithInit(classFile, 0)
}

// ParseGradebookFile parses a gradebook file into a pointer to Gradebook.
func ParseGradebookFile(gradebookFile string) (*Gradebook, error) {
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
	defer c.projectTrustedStudentsToCompatibilityFields()

	return c.loadGradebooks(dir, term, false)
}

// LoadUnscored scans a given directory for *.gradebook files and counts
// unscored assignments for each student by assignment category. The method
// returns an error if there is a problem reading, unmarshaling, or closing a
// file.
func (c *Class) LoadUnscored(dir string, term *Term) error {
	c.initializeStudentMaps(initUnscored)
	defer c.projectTrustedStudentsToCompatibilityFields()

	return c.loadGradebooks(dir, term, true)
}

func (c *Class) trustParsedData() error {
	domain, err := parseTrustedClassDomain(c)
	if err != nil {
		return err
	}

	c.domain = domain
	c.trusted = true
	c.projectTrustedDomainToCompatibilityFields()

	return nil
}

func (c *Class) projectTrustedDomainToCompatibilityFields() {
	domain := c.trustedDomain()
	if domain == nil {
		return
	}

	c.termsByID = cloneTermsByID(domain.termsByID)
	c.assignmentCategories = slices.Clone(domain.assignmentCategories)
	c.categoriesByAssignmentType = maps.Clone(domain.categoriesByType)
	c.labelsByAssignmentCategory = maps.Clone(domain.labelsByCategory)
	c.weightsByAssignmentCategory = maps.Clone(domain.weightsByCategory)
	c.studentsByEmail = cloneStudentsByEmail(domain.studentsByEmail)
}

func (c *Class) projectTrustedStudentsToCompatibilityFields() {
	domain := c.trustedDomain()
	if domain == nil {
		return
	}

	c.studentsByEmail = cloneStudentsByEmail(domain.studentsByEmail)
}

func (c *Class) trustedDomain() *trustedClassDomain {
	if c == nil || !c.trusted || c.domain == nil {
		return nil
	}

	return c.domain
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

		if err := c.loadGradebookFile(gradebook, term, countUnscored); err != nil {
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

func (c *Class) loadGradebookFile(gradebookPath string, term *Term, countUnscored bool) error {
	parsed, err := c.parseGradebookForLoad(gradebookPath)
	if err != nil {
		return err
	}

	// If term != nil, then we are calculating grades for only some of a class.
	// In that case, we want to ignore finals since those grades only matter
	// when we calculate the entire duration of a class.
	if term != nil && parsed.assignmentType == "final" {
		return nil
	}

	for _, record := range parsed.records {
		if countUnscored {
			if record.hasGrade {
				continue
			}
			record.student.unscoredByCategory[parsed.category]++

			continue
		}

		if !record.hasGrade {
			continue
		}

		if !c.trusted {
			if _, ok := record.student.gradesByCategory[parsed.category]; !ok {
				return fmt.Errorf(
					"gradebook: unrecognized assignment category %q for type %q",
					parsed.category,
					parsed.assignmentType,
				)
			}
		}

		record.student.gradesByCategory[parsed.category] = append(
			record.student.gradesByCategory[parsed.category],
			record.grade,
		)
	}

	return nil
}

type parsedAssignmentRecord struct {
	student  *Student
	hasGrade bool
	grade    float64
}

type parsedGradebook struct {
	assignmentType string
	category       string
	records        []parsedAssignmentRecord
}

func (c *Class) parseGradebookForLoad(gradebookPath string) (*parsedGradebook, error) {
	gbData, err := ParseGradebookFile(gradebookPath)
	if err != nil {
		return nil, err
	}

	category, err := c.categoryForAssignmentType(gbData.AssignmentType)
	if err != nil {
		return nil, err
	}

	records := make([]parsedAssignmentRecord, 0, len(gbData.AssignmentRecords))
	for i, ar := range gbData.AssignmentRecords {
		if ar == nil {
			return nil, fmt.Errorf("gradebook: nil assignment record at index %d in %q", i, gradebookPath)
		}

		student, err := c.studentByEmail(ar.Email)
		if err != nil {
			return nil, err
		}

		record := parsedAssignmentRecord{student: student}
		if ar.Grade != nil {
			record.hasGrade = true
			record.grade = *ar.Grade
		}
		records = append(records, record)
	}

	return &parsedGradebook{
		assignmentType: gbData.AssignmentType,
		category:       category,
		records:        records,
	}, nil
}

func (c *Class) studentByEmail(email string) (*Student, error) {
	if domain := c.trustedDomain(); domain != nil {
		student, ok := domain.studentsByEmail[email]
		if !ok {
			return nil, fmt.Errorf("gradebook: no student with email %q", email)
		}
		if student == nil {
			return nil, fmt.Errorf("gradebook: student with email %q is nil", email)
		}

		return student, nil
	}

	student, ok := c.studentsByEmail[email]
	if !ok {
		return nil, fmt.Errorf("gradebook: no student with email %q", email)
	}
	if student == nil {
		return nil, fmt.Errorf("gradebook: student with email %q is nil", email)
	}

	return student, nil
}

func (c *Class) categoryForAssignmentType(assignmentType string) (string, error) {
	if domain := c.trustedDomain(); domain != nil {
		category, ok := domain.categoriesByType[assignmentType]
		if !ok {
			return "", fmt.Errorf("gradebook: unrecognized assignment type %q", assignmentType)
		}

		return category, nil
	}

	category, ok := c.AssignmentCategoryForType(assignmentType)
	if !ok {
		return "", fmt.Errorf("gradebook: unrecognized assignment type %q", assignmentType)
	}

	return category, nil
}
