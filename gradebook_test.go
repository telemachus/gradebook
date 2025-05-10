package gradebook_test

import (
	"gradebook"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var (
	classJSON            = "testdata/class.json"
	classUnequalJSON     = "testdata/wrong.json"
	classInvalidJSON     = "testdata/invalid.json"
	gradebookJSON        = "testdata/quiz-golden-20240319.gradebook"
	gradebookUnequalJSON = "testdata/quiz-wrong-20240319.gradebook"
	gradebookInvalidJSON = "testdata/quiz-invalid-20240319.gradebook"
)

func TestUnmarshalClass(t *testing.T) {
	t.Parallel()

	want := fakeClass()
	got, err := gradebook.UnmarshalClass(classJSON)
	if err != nil {
		t.Fatalf("gradebook.UnmarshalClass(classJSON) = %v; want nil error", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("gradebook.UnmarshalClass(classJSON) mismatch(-want +got):\n%s", diff)
	}
}

func TestUnmarshalClassUnequal(t *testing.T) {
	t.Parallel()

	want := fakeClass()
	got, err := gradebook.UnmarshalClass(classUnequalJSON)
	if err != nil {
		t.Fatalf("gradebook.UnmarshalClass(classUnequalJSON) = %v; want nil error", err)
	}

	if cmp.Equal(want, got) {
		t.Error("gradebook.UnmarshalClass(classUnequalJSON) should differ from the mock class, but it does not")
	}
}

func TestUnmarshalClassInvalid(t *testing.T) {
	t.Parallel()

	_, err := gradebook.UnmarshalClass(classInvalidJSON)
	if err == nil {
		t.Fatal("want error; got nil")
	}
}

func fakeClass() *gradebook.Class {
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
			},
			"mfrede@school.edu": &gradebook.Student{
				FirstName: "Michael",
				LastName:  "Frede",
			},
			"jannas@school.edu": &gradebook.Student{
				FirstName: "Julia",
				LastName:  "Annas",
			},
			"agomezlobo@school.edu": &gradebook.Student{
				FirstName: "Alfonso",
				LastName:  "Gómez-Lobo",
			},
			"gfine@school.edu": &gradebook.Student{
				FirstName: "Gail",
				LastName:  "Fine",
			},
		},
	}
}

func TestUnmarshalGradebook(t *testing.T) {
	t.Parallel()

	want := fakeGradebook()
	got, err := gradebook.UnmarshalGradebook(gradebookJSON)
	if err != nil {
		t.Fatalf("gradebook.UnmarshalGradebook(gradebookJSON) = %v; want nil error", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("gradebook.UnmarshalGradebook(gradebookJSON) mismatch(-want +got):\n%s", diff)
	}
}

func TestUnmarshalGradebookUnequal(t *testing.T) {
	t.Parallel()

	want := fakeGradebook()
	got, err := gradebook.UnmarshalGradebook(gradebookUnequalJSON)
	if err != nil {
		t.Fatalf("gradebook.UnmarshalGradebook(gradebookUnequalJSON) = %v; want nil error", err)
	}

	if cmp.Equal(want, got) {
		t.Error("gradebook.UnmarshalGradebook(gradebookUnequalJSON) should differ from the mock class, but it does not")
	}
}

func TestUnmarshalGradebookInvalid(t *testing.T) {
	t.Parallel()

	_, err := gradebook.UnmarshalGradebook(gradebookInvalidJSON)
	if err == nil {
		t.Fatal("want error; got nil")
	}
}

func fakeGradebook() *gradebook.Gradebook {
	return &gradebook.Gradebook{
		AssignmentDate:        "20240319",
		AssignmentCategory:    "minor",
		AssignmentSubcategory: "quiz",
		Grades: gradebook.Grades{
			&gradebook.Grade{
				Email: "gstriker@school.edu",
				Grade: floatPtr(94.2),
			},
			&gradebook.Grade{
				Email: "mfrede@school.edu",
				Grade: floatPtr(94.0),
			},
			&gradebook.Grade{
				Email: "jannas@school.edu",
				Grade: floatPtr(104),
			},
			&gradebook.Grade{
				Email: "agomezlobo@school.edu",
				Grade: floatPtr(81),
			},
			&gradebook.Grade{
				Email: "gfine@school.edu",
			},
		},
	}
}

func floatPtr(n float64) *float64 {
	return &n
}
