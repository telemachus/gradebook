package gradebook

import (
	"maps"
	"slices"
)

// StudentSpec describes student data for manual class construction.
type StudentSpec struct {
	GradesByCategory   map[string][]float64
	UnscoredByCategory map[string]int
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
}

// StudentsByEmailSpec maps student specs by email.
type StudentsByEmailSpec map[string]*StudentSpec

// ClassSpec describes class data for manual class construction.
type ClassSpec struct {
	TermsByID                   `json:"terms_by_id"`
	LabelsByAssignmentCategory  `json:"labels_by_assignment_category"`
	WeightsByAssignmentCategory `json:"weights_by_assignment_category"`
	CategoriesByAssignmentType  `json:"categories_by_assignment_type"`
	StudentsByEmail             StudentsByEmailSpec `json:"students_by_email"`
	Name                        string              `json:"name"`
	AssignmentCategories        `json:"assignment_categories"`
}

func cloneStudentFromSpec(spec *StudentSpec) *Student {
	if spec == nil {
		return nil
	}

	return &Student{
		gradesByCategory:   cloneGradesByCategory(spec.GradesByCategory),
		unscoredByCategory: maps.Clone(spec.UnscoredByCategory),
		firstName:          spec.FirstName,
		lastName:           spec.LastName,
	}
}

func cloneStudentsByEmailFromSpec(students StudentsByEmailSpec) map[string]*Student {
	if students == nil {
		return nil
	}

	cloned := make(map[string]*Student, len(students))
	for email, spec := range students {
		cloned[email] = cloneStudentFromSpec(spec)
	}

	return cloned
}

func cloneClassFromSpec(spec *ClassSpec) *Class {
	if spec == nil {
		return nil
	}

	return &Class{
		name:                        spec.Name,
		termsByID:                   cloneTermsByID(spec.TermsByID),
		assignmentCategories:        slices.Clone(spec.AssignmentCategories),
		labelsByAssignmentCategory:  maps.Clone(spec.LabelsByAssignmentCategory),
		weightsByAssignmentCategory: maps.Clone(spec.WeightsByAssignmentCategory),
		categoriesByAssignmentType:  maps.Clone(spec.CategoriesByAssignmentType),
		studentsByEmail:             cloneStudentsByEmailFromSpec(spec.StudentsByEmail),
	}
}
