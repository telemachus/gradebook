package gradebook

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestParseClassFile(t *testing.T) {
	t.Parallel()

	want := fakeClass()

	got, err := ParseClassFile(classJSON)
	if err != nil {
		t.Fatalf("ParseClassFile(%q) returned error: %v", classJSON, err)
	}

	if diff := cmp.Diff(
		want,
		got,
		cmpopts.IgnoreFields(Class{}, "trusted", "domain"),
		cmp.AllowUnexported(Class{}, Student{}),
	); diff != "" {
		t.Fatalf("ParseClassFile mismatch (-want +got):\n%s", diff)
	}
}

func TestParseGradebookFile(t *testing.T) {
	t.Parallel()

	want := fakeGradebook()

	got, err := ParseGradebookFile(gradebookJSON)
	if err != nil {
		t.Fatalf("ParseGradebookFile(%q) returned error: %v", gradebookJSON, err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("ParseGradebookFile mismatch (-want +got):\n%s", diff)
	}
}

func TestParseClassFileForGrades(t *testing.T) {
	t.Parallel()

	class, err := ParseClassFileForGrades(classJSON)
	if err != nil {
		t.Fatalf("ParseClassFileForGrades(%q) returned error: %v", classJSON, err)
	}

	student := class.studentsByEmail["gstriker@school.edu"]
	if student == nil {
		t.Fatal("student gstriker@school.edu missing")
	}
	if student.gradesByCategory == nil {
		t.Fatal("GradesByCategory should be initialized")
	}
}

func TestParseClassFileForUnscored(t *testing.T) {
	t.Parallel()

	class, err := ParseClassFileForUnscored(classJSON)
	if err != nil {
		t.Fatalf("ParseClassFileForUnscored(%q) returned error: %v", classJSON, err)
	}

	student := class.studentsByEmail["gstriker@school.edu"]
	if student == nil {
		t.Fatal("student gstriker@school.edu missing")
	}
	if student.unscoredByCategory == nil {
		t.Fatal("UnscoredByCategory should be initialized")
	}
}

func TestClassAccessors(t *testing.T) {
	t.Parallel()

	class, err := ParseClassFile(classJSON)
	if err != nil {
		t.Fatalf("ParseClassFile(%q) returned error: %v", classJSON, err)
	}

	if _, ok := class.TermByID("q1"); !ok {
		t.Fatal("TermByID(q1) should return ok=true")
	}
	if _, ok := class.AssignmentCategoryForType("quiz"); !ok {
		t.Fatal("AssignmentCategoryForType(quiz) should return ok=true")
	}
	if _, ok := class.AssignmentLabelByCategory("major"); !ok {
		t.Fatal("AssignmentLabelByCategory(major) should return ok=true")
	}

	weights := class.AssignmentCategoryWeights()
	weights["major"] = 0
	if class.weightsByAssignmentCategory["major"] == 0 {
		t.Fatal("AssignmentCategoryWeights should return a copy")
	}
}

func TestParseClassFileInvalidClassError(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "class.json")
	data := `{
    "name": "Bad Class",
    "terms_by_id": {
        "q1": {
            "start": "20240101",
            "end": "20240201"
        }
    },
    "assignment_categories": ["major"],
    "labels_by_assignment_category": {
        "major": "Major"
    },
    "weights_by_assignment_category": {
        "major": 100
    },
    "categories_by_assignment_type": {
        "test": "major"
    },
    "students_by_email": {
        "alice.example.com": {
            "first_name": "Alice",
            "last_name": "Zephyr"
        }
    }
}`
	if err := os.WriteFile(file, []byte(data), 0o644); err != nil {
		t.Fatalf("os.WriteFile(%q) failed: %v", file, err)
	}

	_, err := ParseClassFile(file)
	if err == nil {
		t.Fatal("ParseClassFile() should return an error for invalid class data")
	}

	var invalidClassErr *InvalidClassError
	if !errors.As(err, &invalidClassErr) {
		t.Fatalf("ParseClassFile() error type = %T; want *InvalidClassError", err)
	}
}

func TestLoadGradesUsesParsedDomainSnapshot(t *testing.T) {
	t.Parallel()

	class, err := ParseClassFileForGrades(classJSON)
	if err != nil {
		t.Fatalf("ParseClassFileForGrades(%q) returned error: %v", classJSON, err)
	}

	class.categoriesByAssignmentType["quiz"] = "unknown"
	delete(class.studentsByEmail, "gstriker@school.edu")
	class.assignmentCategories = AssignmentCategories{"major"}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() returned error: %v", err)
	}
	validDir := filepath.Join(wd, "testdata", "validgradebooks")

	if err := class.LoadGrades(validDir, nil); err != nil {
		t.Fatalf("LoadGrades(%q, nil) returned error: %v", validDir, err)
	}
}

func TestAccessorsUseParsedDomainSnapshot(t *testing.T) {
	t.Parallel()

	class, err := ParseClassFile(classJSON)
	if err != nil {
		t.Fatalf("ParseClassFile(%q) returned error: %v", classJSON, err)
	}

	class.termsByID["q1"].Start = "20990101"
	class.categoriesByAssignmentType["quiz"] = "major"
	class.labelsByAssignmentCategory["minor"] = "Changed"
	class.weightsByAssignmentCategory["minor"] = 0
	delete(class.categoriesByAssignmentType, "quiz")

	term, ok := class.TermByID("q1")
	if !ok {
		t.Fatal("TermByID(q1) should return ok=true")
	}
	if term.Start != "20200910" {
		t.Fatalf("TermByID(q1).Start = %q; want %q", term.Start, "20200910")
	}

	category, ok := class.AssignmentCategoryForType("quiz")
	if !ok {
		t.Fatal("AssignmentCategoryForType(quiz) should return ok=true")
	}
	if category != "minor" {
		t.Fatalf("AssignmentCategoryForType(quiz) = %q; want %q", category, "minor")
	}

	label, ok := class.AssignmentLabelByCategory("minor")
	if !ok {
		t.Fatal("AssignmentLabelByCategory(minor) should return ok=true")
	}
	if label != "Daily work and quizzes" {
		t.Fatalf("AssignmentLabelByCategory(minor) = %q; want %q", label, "Daily work and quizzes")
	}

	weights := class.AssignmentCategoryWeights()
	if weights["minor"] != 40 {
		t.Fatalf("AssignmentCategoryWeights()[minor] = %d; want %d", weights["minor"], 40)
	}
}

func TestParsedClassTermByIDReturnsCopy(t *testing.T) {
	t.Parallel()

	class, err := ParseClassFile(classJSON)
	if err != nil {
		t.Fatalf("ParseClassFile(%q) returned error: %v", classJSON, err)
	}

	term, ok := class.TermByID("q1")
	if !ok {
		t.Fatal("TermByID(q1) should return ok=true")
	}
	term.Start = "20990101"

	again, ok := class.TermByID("q1")
	if !ok {
		t.Fatal("TermByID(q1) should return ok=true on second call")
	}
	if again.Start != "20200910" {
		t.Fatalf("TermByID(q1).Start = %q; want %q", again.Start, "20200910")
	}
}

func TestSortingUsesParsedDomainSnapshot(t *testing.T) {
	t.Parallel()

	class, err := ParseClassFile(classJSON)
	if err != nil {
		t.Fatalf("ParseClassFile(%q) returned error: %v", classJSON, err)
	}

	class.assignmentCategories = AssignmentCategories{"major"}
	class.labelsByAssignmentCategory["major"] = "zzz"
	class.studentsByEmail = StudentsByEmail{}

	categories := class.AssignmentCategoriesSortedByLabel()
	if len(categories) != 3 {
		t.Fatalf("len(AssignmentCategoriesSortedByLabel()) = %d; want %d", len(categories), 3)
	}

	emails := class.EmailsSortedByStudentName()
	if len(emails) != 5 {
		t.Fatalf("len(EmailsSortedByStudentName()) = %d; want %d", len(emails), 5)
	}

	students := class.StudentsSortedByName()
	if len(students) != 5 {
		t.Fatalf("len(StudentsSortedByName()) = %d; want %d", len(students), 5)
	}
}

func TestParsedClassStudentMutationsDoNotChangeDomainSorting(t *testing.T) {
	t.Parallel()

	class, err := ParseClassFile(classJSON)
	if err != nil {
		t.Fatalf("ParseClassFile(%q) returned error: %v", classJSON, err)
	}

	before := class.EmailsSortedByStudentName()
	class.studentsByEmail["gstriker@school.edu"].lastName = "Aardvark"
	after := class.EmailsSortedByStudentName()

	if diff := cmp.Diff(before, after); diff != "" {
		t.Fatalf("EmailsSortedByStudentName changed after external student mutation (-before +after):\n%s", diff)
	}
}

func TestParsedClassProjectsCanonicalStudentsAfterLoadGrades(t *testing.T) {
	t.Parallel()

	class, err := ParseClassFileForGrades(classJSON)
	if err != nil {
		t.Fatalf("ParseClassFileForGrades(%q) returned error: %v", classJSON, err)
	}

	replacement := &Student{
		firstName: "Replacement",
		lastName:  "Student",
		gradesByCategory: map[string][]float64{
			"major": {},
			"minor": {},
			"cp":    {},
		},
	}
	class.studentsByEmail["gstriker@school.edu"] = replacement

	dir := t.TempDir()
	gradebookData := `{
    "assignment_date": "20240319",
    "assignment_name": "projection-check",
    "assignment_type": "quiz",
    "assignment_category": "minor",
    "assignment_records": [
        {
            "email": "gstriker@school.edu",
            "grade": 88
        }
    ]
}`
	file := filepath.Join(dir, "quiz-projection-check-20240319.gradebook")
	if err := os.WriteFile(file, []byte(gradebookData), 0o644); err != nil {
		t.Fatalf("os.WriteFile(%q) failed: %v", file, err)
	}

	if err := class.LoadGrades(dir, nil); err != nil {
		t.Fatalf("LoadGrades(%q, nil) returned error: %v", dir, err)
	}

	student := class.studentsByEmail["gstriker@school.edu"]
	if student == replacement {
		t.Fatal("LoadGrades should re-project trusted canonical students")
	}
	if got := len(student.gradesByCategory["minor"]); got != 1 {
		t.Fatalf("len(GradesByCategory[minor]) = %d; want %d", got, 1)
	}
	if got := student.gradesByCategory["minor"][0]; got != 88 {
		t.Fatalf("GradesByCategory[minor][0] = %v; want %v", got, 88)
	}
	if got := len(replacement.gradesByCategory["minor"]); got != 0 {
		t.Fatalf("replacement student was mutated; got %d grades", got)
	}
}

func TestParsedClassStudentsSortedByNameReturnsCopies(t *testing.T) {
	t.Parallel()

	class, err := ParseClassFile(classJSON)
	if err != nil {
		t.Fatalf("ParseClassFile(%q) returned error: %v", classJSON, err)
	}

	students := class.StudentsSortedByName()
	students[0].lastName = "Mutated"

	again := class.StudentsSortedByName()
	if again[0].lastName == "Mutated" {
		t.Fatal("StudentsSortedByName should return student copies for parsed classes")
	}
}
