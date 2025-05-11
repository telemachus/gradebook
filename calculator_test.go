package gradebook_test

import (
	"gradebook"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUnmarshalCalcClass(t *testing.T) {
	t.Parallel()

	want := fakeCalcClass()
	got, err := gradebook.UnmarshalCalcClass(classJSON)
	if err != nil {
		t.Fatalf("gradebook.UnmarshalClass(classJSON) = %v; want nil error", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("gradebook.UnmarshalClass(classJSON) mismatch(-want +got):\n%s", diff)
	}
}

func fakeGradesMap() map[string][]*float64 {
	gradesByCategories := make(map[string][]*float64, 3)
	gradesByCategories["major"] = make([]*float64, 0, 50)
	gradesByCategories["minor"] = make([]*float64, 0, 50)
	gradesByCategories["cp"] = make([]*float64, 0, 50)

	return gradesByCategories
}

func fakeCalcClass() *gradebook.Class {
	return &gradebook.Class{
		Name: "Lucretius",
		TermsByID: map[string]*gradebook.Term{
			"q1": {
				Start: "20200910",
				End:   "20201103",
			},
			"q2": {
				Start: "20201108",
				End:   "20210114",
			},
			"q3": {
				Start: "20210124",
				End:   "20210311",
			},
			"q4": {
				Start: "20210328",
				End:   "20210609",
			},
			"s1": {
				Start: "20200910",
				End:   "20210114",
			},
			"s2": {
				Start: "20210124",
				End:   "20210609",
			},
		},
		AssignmentTypes: gradebook.AssignmentTypes{"major", "minor", "cp"},
		LabelsByAssignmentType: gradebook.LabelsByAssignmentType{
			"major": "Major assessments",
			"minor": "Daily work and quizzes",
			"cp":    "Class participation",
		},
		WeightsByAssignmentType: gradebook.WeightsByAssignmentType{
			"major": 30,
			"minor": 40,
			"cp":    30,
		},
		CategoriesByAssignmentType: gradebook.CategoriesByAssignmentType{
			"test":    "major",
			"project": "major",
			"essay":   "major",
			"quiz":    "minor",
			"hw":      "minor",
			"cp":      "cp",
		},
		StudentsByEmail: gradebook.StudentsByEmail{
			"gstriker@school.edu": &gradebook.Student{
				FirstName:    "Gisela",
				LastName:     "Striker",
				GradesByType: fakeGradesMap(),
			},
			"mfrede@school.edu": &gradebook.Student{
				FirstName:    "Michael",
				LastName:     "Frede",
				GradesByType: fakeGradesMap(),
			},
			"jannas@school.edu": &gradebook.Student{
				FirstName:    "Julia",
				LastName:     "Annas",
				GradesByType: fakeGradesMap(),
			},
			"agomezlobo@school.edu": &gradebook.Student{
				FirstName:    "Alfonso",
				LastName:     "Gómez-Lobo",
				GradesByType: fakeGradesMap(),
			},
			"gfine@school.edu": &gradebook.Student{
				FirstName:    "Gail",
				LastName:     "Fine",
				GradesByType: fakeGradesMap(),
			},
		},
	}
}
