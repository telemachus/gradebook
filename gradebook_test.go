package gradebook_test

import (
	"gradebook"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var (
	identicalJSON = "testdata/identical.json"
	differentJSON = "testdata/different.json"
	invalidJSON   = "testdata/invalid.json"
)

func TestUnmarshalClassIdenticalValid(t *testing.T) {
	t.Parallel()

	want := testClass()

	got, err := gradebook.UnmarshalClass(identicalJSON)
	if err != nil {
		t.Fatalf("gradebook.UnmarshalClass(identicalJSON) = %v; want nil error", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("gradebook.UnmarshalClass(identicalJSON) mismatch(-want +got):\n%s", diff)
	}
}

func TestUnmarshalClassDifferentValid(t *testing.T) {
	t.Parallel()

	want := testClass()

	got, err := gradebook.UnmarshalClass(differentJSON)
	if err != nil {
		t.Fatalf("gradebook.UnmarshalClass(differentJSON) = %v; want nil error", err)
	}

	if cmp.Equal(want, got) {
		t.Error("gradebook.UnmarshalClass(differentJSON) should differ from the mock class, but it does not")
	}
}

func TestUnmarshalClassInvalid(t *testing.T) {
	t.Parallel()

	_, err := gradebook.UnmarshalClass(invalidJSON)

	if err == nil {
		t.Fatal("want error; got nil")
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
