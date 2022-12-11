package gradebook_test

import (
	"gradebook"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var identicalJSON = "testdata/identical.json"
var differentJSON = "testdata/different.json"
var invalidJSON = "testdata/invalid.json"

func TestLoadClassIdenticalValid(t *testing.T) {
	t.Parallel()

	expected := testClass()
	actual, err := gradebook.LoadClass(identicalJSON)

	if err != nil {
		t.Fatalf("expected nil error; actual %v", err)
	}

	if !cmp.Equal(expected, actual) {
		t.Error(cmp.Diff(expected, actual))
	}
}

func TestLoadClassDifferentValid(t *testing.T) {
	t.Parallel()

	expected := testClass()
	actual, err := gradebook.LoadClass(differentJSON)

	if err != nil {
		t.Fatalf("expected nil error; actual %v", err)
	}

	if cmp.Equal(expected, actual) {
		t.Error("expected different classes; actual identical classes")
	}
}

func TestLoadClassInvalid(t *testing.T) {
	t.Parallel()

	_, err := gradebook.LoadClass(invalidJSON)

	if err == nil {
		t.Fatal("expected error; actual nil")
	}
}

func testClass() *gradebook.Class {
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
		CategoriesPretty: gradebook.CategoriesPretty{
			"major": "Major assessments",
			"minor": "Daily work and quizzes",
			"cp":    "Class participation",
		},
		CategoryWeights: gradebook.CategoryWeights{
			"major": 30,
			"minor": 40,
			"cp":    30,
		},
		TypesToCategories: gradebook.TypesToCategories{
			"test":    "major",
			"project": "major",
			"essay":   "major",
			"quiz":    "minor",
			"hw":      "minor",
			"cp":      "cp",
		},
		Students: gradebook.Students{
			"gstriker@school.edu": {
				FirstName: "Gisela",
				LastName:  "Striker",
			},
			"mfrede@school.edu": {
				FirstName: "Michael",
				LastName:  "Frede",
			},
			"jannas@school.edu": {
				FirstName: "Julia",
				LastName:  "Annas",
			},
			"agomezlobo@school.edu": {
				FirstName: "Alfonso",
				LastName:  "Gómez-Lobo",
			},
			"gfine@school.edu": {
				FirstName: "Gail",
				LastName:  "Fine",
			},
		},
	}
}
