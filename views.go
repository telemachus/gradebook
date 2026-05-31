package gradebook

import (
	"maps"
	"slices"
)

// Name returns the class name.
func (c *Class) Name() string {
	if c == nil {
		return ""
	}

	return c.name
}

// TermsByID returns a copy of terms by ID.
func (c *Class) TermsByID() TermsByID {
	if c == nil {
		return TermsByID{}
	}

	return TermsByID(cloneTermsByID(c.termsForLookup()))
}

// AssignmentCategories returns a copy of assignment categories.
func (c *Class) AssignmentCategories() AssignmentCategories {
	if c == nil {
		return AssignmentCategories{}
	}
	if domain := c.trustedDomain(); domain != nil {
		return slices.Clone(domain.assignmentCategories)
	}

	return slices.Clone(c.assignmentCategories)
}

// LabelsByAssignmentCategory returns a copy of labels by assignment category.
func (c *Class) LabelsByAssignmentCategory() LabelsByAssignmentCategory {
	if c == nil {
		return LabelsByAssignmentCategory{}
	}

	return maps.Clone(c.labelsByCategoryForLookup())
}

// CategoriesByAssignmentType returns a copy of categories by assignment type.
func (c *Class) CategoriesByAssignmentType() CategoriesByAssignmentType {
	if c == nil {
		return CategoriesByAssignmentType{}
	}

	return maps.Clone(c.categoriesByTypeForLookup())
}

// StudentsByEmail returns copies of students by email.
func (c *Class) StudentsByEmail() StudentsByEmail {
	if c == nil {
		return StudentsByEmail{}
	}
	if domain := c.trustedDomain(); domain != nil {
		return StudentsByEmail(cloneStudentsByEmail(domain.studentsByEmail))
	}

	return StudentsByEmail(cloneStudentsByEmail(c.studentsByEmail))
}

// FirstName returns the student's first name.
func (s *Student) FirstName() string {
	if s == nil {
		return ""
	}

	return s.firstName
}

// LastName returns the student's last name.
func (s *Student) LastName() string {
	if s == nil {
		return ""
	}

	return s.lastName
}

// GradesByCategory returns a copy of grades by category.
func (s *Student) GradesByCategory() map[string][]float64 {
	if s == nil {
		return map[string][]float64{}
	}

	return cloneGradesByCategory(s.gradesByCategory)
}

// UnscoredByCategory returns a copy of unscored counts by category.
func (s *Student) UnscoredByCategory() map[string]int {
	if s == nil {
		return map[string]int{}
	}

	return maps.Clone(s.unscoredByCategory)
}

// UnscoredCountByCategory returns unscored count for a category.
func (s *Student) UnscoredCountByCategory(category string) int {
	if s == nil {
		return 0
	}

	return s.unscoredByCategory[category]
}
