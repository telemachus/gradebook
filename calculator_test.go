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
	gMap := make(map[string][]*float64, 3)
	gMap["major"] = make([]*float64, 0, 50)
	gMap["minor"] = make([]*float64, 0, 50)
	gMap["cp"] = make([]*float64, 0, 50)
	return gMap
}

func fakeCalcClass() *gradebook.Class {
	return &gradebook.Class{
		Name: "Lucretius",
		Terms: map[string]*gradebook.Term{
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
		Categories: gradebook.Categories{"major", "minor", "cp"},
		PrettyCategories: gradebook.PrettyCategories{
			"major": "Major assessments",
			"minor": "Daily work and quizzes",
			"cp":    "Class participation",
		},
		Weights: gradebook.Weights{
			"major": 30,
			"minor": 40,
			"cp":    30,
		},
		Subcategories: gradebook.Subcategories{
			"test":    "major",
			"project": "major",
			"essay":   "major",
			"quiz":    "minor",
			"hw":      "minor",
			"cp":      "cp",
		},
		Students: gradebook.Students{
			"gstriker@school.edu": &gradebook.Student{
				FirstName: "Gisela",
				LastName:  "Striker",
				Grades:    fakeGradesMap(),
			},
			"mfrede@school.edu": &gradebook.Student{
				FirstName: "Michael",
				LastName:  "Frede",
				Grades:    fakeGradesMap(),
			},
			"jannas@school.edu": &gradebook.Student{
				FirstName: "Julia",
				LastName:  "Annas",
				Grades:    fakeGradesMap(),
			},
			"agomezlobo@school.edu": &gradebook.Student{
				FirstName: "Alfonso",
				LastName:  "Gómez-Lobo",
				Grades:    fakeGradesMap(),
			},
			"gfine@school.edu": &gradebook.Student{
				FirstName: "Gail",
				LastName:  "Fine",
				Grades:    fakeGradesMap(),
			},
		},
	}
}
