package gradebook_test

import (
	"gradebook"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Let's build a class.
var q1 = gradebook.Term{
	Start: "20200910",
	End:   "20201103",
}
var q2 = gradebook.Term{
	Start: "20201108",
	End:   "20210114",
}
var q3 = gradebook.Term{
	Start: "20210124",
	End:   "20210311",
}
var q4 = gradebook.Term{
	Start: "20210328",
	End:   "20210609",
}
var s1 = gradebook.Term{
	Start: q1.Start,
	End:   q2.End,
}
var s2 = gradebook.Term{
	Start: q3.Start,
	End:   q4.End,
}
var terms = gradebook.Terms{
	"q1": &q1,
	"q2": &q2,
	"q3": &q3,
	"q4": &q4,
	"s1": &s1,
	"s2": &s2,
}

var categories = gradebook.Categories{"major", "minor", "cp"}
var categoriesPretty = gradebook.CategoriesPretty{
	categories[0]: "Major assessments",
	categories[1]: "Daily work and quizzes",
	categories[2]: "Class participation",
}
var categoryWeights = gradebook.CategoryWeights{
	categories[0]: 30,
	categories[1]: 40,
	categories[2]: 30,
}
var typesToCategories = gradebook.TypesToCategories{
	"test":    categories[0],
	"project": categories[0],
	"essay":   categories[0],
	"quiz":    categories[1],
	"hw":      categories[1],
	"cp":      categories[2],
}

var gstriker = gradebook.Student{
	FirstName: "Gisela",
	LastName:  "Striker",
}
var mfrede = gradebook.Student{
	FirstName: "Michael",
	LastName:  "Frede",
}
var jannas = gradebook.Student{
	FirstName: "Julia",
	LastName:  "Annas",
}
var agomezlobo = gradebook.Student{
	FirstName: "Alfonso",
	LastName:  "Gómez-Lobo",
}
var gfine = gradebook.Student{
	FirstName: "Gail",
	LastName:  "Fine",
}
var students = gradebook.Students{
	"gstriker@school.edu":   &gstriker,
	"mfrede@school.edu":     &mfrede,
	"jannas@school.edu":     &jannas,
	"agomezlobo@school.edu": &agomezlobo,
	"gfine@school.edu":      &gfine,
}

var class = gradebook.Class{
	Name:              "Lucretius",
	Terms:             terms,
	Categories:        categories,
	CategoriesPretty:  categoriesPretty,
	CategoryWeights:   categoryWeights,
	TypesToCategories: typesToCategories,
	Students:          students,
}

var identicalJSON = "testdata/identical.json"
var differentJSON = "testdata/different.json"
var invalidJSON = "testdata/invalid.json"

func TestLoadClassIdenticalValid(t *testing.T) {
	t.Parallel()

	expected := &class
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

	expected := &class
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
