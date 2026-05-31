package gradebook

import (
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var (
	classJSON            = "testdata/class.json"
	classUnequalJSON     = "testdata/wrong.json"
	classInvalidJSON     = "testdata/invalid.json"
	gradebookJSON        = "testdata/quiz-golden-20240319.gradebook"
	gradebookUnequalJSON = "testdata/quiz-wrong-20240319.gradebook"
	gradebookInvalidJSON = "testdata/quiz-invalid-20240319.gradebook"
)

func TestParseClassFileEqualMock(t *testing.T) {
	t.Parallel()

	want := fakeClass()
	got, err := ParseClassFile(classJSON)
	if err != nil {
		t.Fatalf("ParseClassFile(classJSON) = %v; want nil error", err)
	}

	if diff := cmp.Diff(
		want,
		got,
		cmpopts.IgnoreFields(Class{}, "trusted", "domain"),
		cmp.AllowUnexported(Class{}, Student{}),
	); diff != "" {
		t.Errorf("ParseClassFile(classJSON) mismatch(-want +got):\n%s", diff)
	}
}

func TestParseClassFileUnequal(t *testing.T) {
	t.Parallel()

	want := fakeClass()
	got, err := ParseClassFile(classUnequalJSON)
	if err != nil {
		t.Fatalf("ParseClassFile(classUnequalJSON) = %v; want nil error", err)
	}

	if cmp.Equal(
		want,
		got,
		cmpopts.IgnoreFields(Class{}, "trusted", "domain"),
		cmp.AllowUnexported(Class{}, Student{}),
	) {
		t.Error("ParseClassFile(classUnequalJSON) should differ from the mock class, but it does not")
	}
}

func TestParseClassFileInvalid(t *testing.T) {
	t.Parallel()

	_, err := ParseClassFile(classInvalidJSON)
	if err == nil {
		t.Fatal("want error; got nil")
	}
}

func TestParseGradebookFileEqualMock(t *testing.T) {
	t.Parallel()

	want := fakeGradebook()
	got, err := ParseGradebookFile(gradebookJSON)
	if err != nil {
		t.Fatalf("ParseGradebookFile(gradebookJSON) = %v; want nil error", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("ParseGradebookFile(gradebookJSON) mismatch(-want +got):\n%s", diff)
	}
}

func TestParseGradebookFileUnequal(t *testing.T) {
	t.Parallel()

	want := fakeGradebook()
	got, err := ParseGradebookFile(gradebookUnequalJSON)
	if err != nil {
		t.Fatalf("ParseGradebookFile(gradebookUnequalJSON) = %v; want nil error", err)
	}

	if cmp.Equal(want, got) {
		t.Error("ParseGradebookFile(gradebookUnequalJSON) should differ from the mock class, but it does not")
	}
}

func TestParseGradebookFileInvalid(t *testing.T) {
	t.Parallel()

	_, err := ParseGradebookFile(gradebookInvalidJSON)
	if err == nil {
		t.Fatalf("ParseGradebookFile(%q) should error; got nil", gradebookInvalidJSON)
	}
}

func TestLoadGradesValid(t *testing.T) {
	t.Parallel()

	pseudoClass := fakeCalcClass()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal("problem with os.Getwd(): cannot continue")
	}
	testDir := filepath.Join(wd, "testdata", "validgradebooks")

	err = pseudoClass.LoadGrades(testDir, nil)
	if err != nil {
		t.Fatalf("want no error from LoadGrades(%q, nil); got %v", testDir, err)
	}
}

func TestLoadGradesInvalid(t *testing.T) {
	t.Parallel()

	pseudoClass := fakeCalcClass()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal("problem with os.Getwd(): cannot continue")
	}
	testDir := filepath.Join(wd, "testdata", "invalidgradebook")

	err = pseudoClass.LoadGrades(testDir, nil)
	if err == nil {
		t.Fatalf("want error from LoadGrades(%q, nil); got nil", testDir)
	}
}

func TestLoadGradesUnknownType(t *testing.T) {
	t.Parallel()

	pseudoClass := fakeCalcClass()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal("problem with os.Getwd(): cannot continue")
	}
	testDir := filepath.Join(wd, "testdata", "unknowntypegradebook")

	err = pseudoClass.LoadGrades(testDir, nil)
	if err == nil {
		t.Fatalf("want error from LoadGrades(%q, nil); got nil", testDir)
	}
}

func TestLoadGradesNonexistentDirectory(t *testing.T) {
	t.Parallel()

	pseudoClass := fakeCalcClass()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal("problem with os.Getwd(): cannot continue")
	}
	testDir := filepath.Join(wd, "testdata", "nonexistent")

	err = pseudoClass.LoadGrades(testDir, nil)
	if err == nil {
		t.Fatalf("want error from LoadGrades(%q, nil); got nil", testDir)
	}
}

func fakeCalcClass() *Class {
	return &Class{
		name:                        "Lucretius",
		termsByID:                   fakeTermsByID(),
		assignmentCategories:        fakeAssignmentCategories(),
		labelsByAssignmentCategory:  fakeLabelsByAssignmentCategory(),
		weightsByAssignmentCategory: fakeWeightsByAssignmentCategory(),
		categoriesByAssignmentType:  fakeCategoriesByAssignmentType(),
		studentsByEmail:             fakeCalcStudentsByEmail(),
	}
}

func fakeTermsByID() map[string]*Term {
	return map[string]*Term{
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
	}
}

func fakeAssignmentCategories() AssignmentCategories {
	return AssignmentCategories{"major", "minor", "cp"}
}

func fakeLabelsByAssignmentCategory() LabelsByAssignmentCategory {
	return LabelsByAssignmentCategory{
		"major": "Major assessments",
		"minor": "Daily work and quizzes",
		"cp":    "Class participation",
	}
}

func fakeWeightsByAssignmentCategory() WeightsByAssignmentCategory {
	return WeightsByAssignmentCategory{
		"major": 30,
		"minor": 40,
		"cp":    30,
	}
}

func fakeCategoriesByAssignmentType() CategoriesByAssignmentType {
	return CategoriesByAssignmentType{
		"test":    "major",
		"project": "major",
		"essay":   "major",
		"quiz":    "minor",
		"hw":      "minor",
		"cp":      "cp",
	}
}

func fakeStudentsByEmail() StudentsByEmail {
	return StudentsByEmail{
		"gstriker@school.edu": &Student{
			firstName: "Gisela",
			lastName:  "Striker",
		},
		"mfrede@school.edu": &Student{
			firstName: "Michael",
			lastName:  "Frede",
		},
		"jannas@school.edu": &Student{
			firstName: "Julia",
			lastName:  "Annas",
		},
		"agomezlobo@school.edu": &Student{
			firstName: "Alfonso",
			lastName:  "Gómez-Lobo",
		},
		"gfine@school.edu": &Student{
			firstName: "Gail",
			lastName:  "Fine",
		},
	}
}

func fakeCalcStudentsByEmail() StudentsByEmail {
	students := fakeStudentsByEmail()
	for _, student := range students {
		student.gradesByCategory = fakeGradesMap()
	}

	return students
}

func fakeGradebook() *Gradebook {
	return &Gradebook{
		AssignmentCategory: "minor",
		AssignmentDate:     "20240319",
		AssignmentName:     "golden",
		AssignmentType:     "quiz",
		AssignmentRecords: AssignmentRecords{
			&AssignmentRecord{
				Email: "gstriker@school.edu",
				Grade: floatPtr(94.2),
			},
			&AssignmentRecord{
				Email: "mfrede@school.edu",
				Grade: floatPtr(94.0),
			},
			&AssignmentRecord{
				Email: "jannas@school.edu",
				Grade: floatPtr(104),
			},
			&AssignmentRecord{
				Email: "agomezlobo@school.edu",
				Grade: floatPtr(81),
			},
			&AssignmentRecord{
				Email: "gfine@school.edu",
			},
		},
	}
}

func fakeClass() *Class {
	return &Class{
		name:                        "Lucretius",
		termsByID:                   fakeTermsByID(),
		assignmentCategories:        fakeAssignmentCategories(),
		labelsByAssignmentCategory:  fakeLabelsByAssignmentCategory(),
		weightsByAssignmentCategory: fakeWeightsByAssignmentCategory(),
		categoriesByAssignmentType:  fakeCategoriesByAssignmentType(),
		studentsByEmail:             fakeStudentsByEmail(),
	}
}

func TestStudentAverage(t *testing.T) {
	t.Parallel()

	student, err := NewStudent("Michael", "Frede")
	if err != nil {
		t.Fatalf("NewStudent() returned error: %v", err)
	}

	student.gradesByCategory = map[string][]float64{
		"major": make([]float64, 0),
	}

	result := student.Average("major")
	if result.Valid {
		t.Error("Average(major) with no grades should return Valid=false")
	}

	grades := []float64{85, 90, 95}
	student.gradesByCategory["major"] = append(student.gradesByCategory["major"], grades...)

	result = student.Average("major")
	if !result.Valid {
		t.Error("Average(major) with grades should return Valid=true")
	}

	expectedAvg := 90.0
	if !floatEqual(result.Value, expectedAvg, 0.001) {
		t.Errorf("Average(major) = %f; want %f", result.Value, expectedAvg)
	}
}

func TestStudentAverageInvalidCategory(t *testing.T) {
	t.Parallel()

	student, err := NewStudent("Michael", "Frede")
	if err != nil {
		t.Fatalf("NewStudent() returned error: %v", err)
	}

	student.gradesByCategory = map[string][]float64{
		"major": make([]float64, 0),
	}
}

func TestStudentTotalAverage(t *testing.T) {
	t.Parallel()

	student, err := NewStudent("Michael", "Frede")
	if err != nil {
		t.Fatalf("NewStudent() returned error: %v", err)
	}

	student.gradesByCategory = map[string][]float64{
		"major": make([]float64, 0),
		"minor": make([]float64, 0),
		"cp":    make([]float64, 0),
	}

	weights := WeightsByAssignmentCategory{
		"major": 50,
		"minor": 30,
		"cp":    20,
	}

	result := student.TotalAverage(weights)
	if result.Valid {
		t.Error("TotalAverage() with no grades should return Valid=false")
	}

	student.gradesByCategory["major"] = append(student.gradesByCategory["major"], 90)
	student.gradesByCategory["minor"] = append(student.gradesByCategory["minor"], 90)
	student.gradesByCategory["cp"] = append(student.gradesByCategory["cp"], 90)

	result = student.TotalAverage(weights)
	if !result.Valid {
		t.Error("TotalAverage() with grades should return Valid=true")
	}

	expectedAvg := 90.0
	if !floatEqual(result.Value, expectedAvg, 0.001) {
		t.Errorf("TotalAverage() with equal grades = %f; want %f", result.Value, expectedAvg)
	}

	student.gradesByCategory = map[string][]float64{
		"major": {94},
		"minor": {82},
		"cp":    {75},
	}

	result = student.TotalAverage(weights)
	if !result.Valid {
		t.Error("TotalAverage() with grades should return Valid=true")
	}

	expectedAvg = 86.6
	if !floatEqual(result.Value, expectedAvg, 0.1) {
		t.Errorf("TotalAverage() with different grades = %f; want %f", result.Value, expectedAvg)
	}
}

func TestStudentTotalAveragePartialGrades(t *testing.T) {
	t.Parallel()

	student, err := NewStudent("Michael", "Frede")
	if err != nil {
		t.Fatalf("NewStudent() returned error: %v", err)
	}

	student.gradesByCategory = map[string][]float64{
		"major": {90},
		"minor": {80},
		"cp":    make([]float64, 0),
	}

	weights := map[string]int{
		"major": 50,
		"minor": 30,
		"cp":    20,
	}

	result := student.TotalAverage(weights)
	if !result.Valid {
		t.Error("TotalAverage() with partial grades should return Valid=true")
	}

	expectedAvg := 86.25
	if !floatEqual(result.Value, expectedAvg, 0.01) {
		t.Errorf("TotalAverage() with partial grades = %f; want %f", result.Value, expectedAvg)
	}
}

func TestStudentAverageMultipleGrades(t *testing.T) {
	t.Parallel()

	student, err := NewStudent("Michael", "Frede")
	if err != nil {
		t.Fatalf("NewStudent() returned error: %v", err)
	}

	student.gradesByCategory = map[string][]float64{
		"major": {88, 92, 85, 95},
	}

	result := student.Average("major")
	if !result.Valid {
		t.Error("Average(major) with grades should return Valid=true")
	}

	expectedAvg := 90.0
	if !floatEqual(result.Value, expectedAvg, 0.001) {
		t.Errorf("Average(major) with multiple grades = %f; want %f", result.Value, expectedAvg)
	}
}

func floatPtr(n float64) *float64 {
	return &n
}

func floatEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) < tolerance
}
