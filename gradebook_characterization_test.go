package gradebook

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadGradesSkipsNullGrades(t *testing.T) {
	t.Parallel()

	class := fakeCalcClass()
	dir := t.TempDir()

	gradebookData := `{
    "assignment_date": "20240319",
    "assignment_name": "null-check",
    "assignment_type": "quiz",
    "assignment_category": "minor",
    "assignment_records": [
        {
            "email": "gstriker@school.edu",
            "grade": 92
        },
        {
            "email": "gfine@school.edu",
            "grade": null
        }
    ]
}`

	if err := os.WriteFile(filepath.Join(dir, "quiz-null-check-20240319.gradebook"), []byte(gradebookData), 0o644); err != nil {
		t.Fatalf("failed writing gradebook fixture: %v", err)
	}

	if err := class.LoadGrades(dir, nil); err != nil {
		t.Fatalf("LoadGrades(%q, nil) returned error: %v", dir, err)
	}

	if got := len(class.studentsByEmail["gstriker@school.edu"].gradesByCategory["minor"]); got != 1 {
		t.Fatalf("expected one non-null grade for gstriker@school.edu, got %d", got)
	}
	if got := class.studentsByEmail["gstriker@school.edu"].gradesByCategory["minor"][0]; got != 92 {
		t.Fatalf("expected gstriker@school.edu minor grade to be 92, got %v", got)
	}

	if got := len(class.studentsByEmail["gfine@school.edu"].gradesByCategory["minor"]); got != 0 {
		t.Fatalf("expected null grade to be skipped for gfine@school.edu, got %d loaded grades", got)
	}
}

func TestLoadGradesUnknownStudentError(t *testing.T) {
	t.Parallel()

	class := fakeCalcClass()
	dir := t.TempDir()

	gradebookData := `{
    "assignment_date": "20240319",
    "assignment_name": "unknown-student",
    "assignment_type": "quiz",
    "assignment_category": "minor",
    "assignment_records": [
        {
            "email": "nobody@school.edu",
            "grade": 88
        }
    ]
}`

	if err := os.WriteFile(filepath.Join(dir, "quiz-unknown-student-20240319.gradebook"), []byte(gradebookData), 0o644); err != nil {
		t.Fatalf("failed writing gradebook fixture: %v", err)
	}

	err := class.LoadGrades(dir, nil)
	if err == nil {
		t.Fatal("LoadGrades should return error for unknown student email; got nil")
	}
	if !strings.Contains(err.Error(), `no student with email "nobody@school.edu"`) {
		t.Fatalf("unexpected error for unknown student: %v", err)
	}
}

func TestLoadGradesUsesTypeToCategoryMapping(t *testing.T) {
	t.Parallel()

	class := fakeCalcClass()
	dir := t.TempDir()

	gradebookData := `{
    "assignment_date": "20240319",
    "assignment_name": "type-mapping",
    "assignment_type": "quiz",
    "assignment_category": "major",
    "assignment_records": [
        {
            "email": "gstriker@school.edu",
            "grade": 83
        }
    ]
}`

	if err := os.WriteFile(filepath.Join(dir, "quiz-type-mapping-20240319.gradebook"), []byte(gradebookData), 0o644); err != nil {
		t.Fatalf("failed writing gradebook fixture: %v", err)
	}

	if err := class.LoadGrades(dir, nil); err != nil {
		t.Fatalf("LoadGrades(%q, nil) returned error: %v", dir, err)
	}

	if got := len(class.studentsByEmail["gstriker@school.edu"].gradesByCategory["minor"]); got != 1 {
		t.Fatalf("expected grade to be loaded into \"minor\" via assignment_type mapping, got %d", got)
	}
	if got := len(class.studentsByEmail["gstriker@school.edu"].gradesByCategory["major"]); got != 0 {
		t.Fatalf("expected no grade in \"major\" category, got %d", got)
	}
}

func TestLoadGradesWithParseClassFileInitializesStudentMaps(t *testing.T) {
	t.Parallel()

	class, err := ParseClassFile(classJSON)
	if err != nil {
		t.Fatalf("ParseClassFile(%q) returned error: %v", classJSON, err)
	}

	if err := class.LoadGrades("testdata/validgradebooks", nil); err != nil {
		t.Fatalf("LoadGrades should initialize grade maps on demand: %v", err)
	}

	if got := len(class.studentsByEmail["gstriker@school.edu"].gradesByCategory["minor"]); got == 0 {
		t.Fatal("expected grades to load when class is created with ParseClassFile")
	}
}

func TestLoadUnscoredWithParseClassFileInitializesStudentMaps(t *testing.T) {
	t.Parallel()

	class, err := ParseClassFile(classJSON)
	if err != nil {
		t.Fatalf("ParseClassFile(%q) returned error: %v", classJSON, err)
	}

	if err := class.LoadUnscored("testdata/validgradebooks", nil); err != nil {
		t.Fatalf("LoadUnscored should initialize maps on demand: %v", err)
	}

	want := map[string]map[string]int{
		"gstriker@school.edu": {
			"major": 0,
			"minor": 0,
			"cp":    0,
		},
		"mfrede@school.edu": {
			"major": 0,
			"minor": 0,
			"cp":    0,
		},
		"jannas@school.edu": {
			"major": 0,
			"minor": 0,
			"cp":    0,
		},
		"agomezlobo@school.edu": {
			"major": 0,
			"minor": 0,
			"cp":    0,
		},
		"gfine@school.edu": {
			"major": 0,
			"minor": 1,
			"cp":    0,
		},
	}

	for email, wantByCategory := range want {
		student := class.studentsByEmail[email]
		if student == nil {
			t.Fatalf("student %q not found", email)
		}

		for category, wantCount := range wantByCategory {
			gotCount := student.unscoredByCategory[category]
			if gotCount != wantCount {
				t.Fatalf(
					"UnscoredByCategory[%q][%q] = %d; want %d",
					email,
					category,
					gotCount,
					wantCount,
				)
			}
		}
	}
}

func TestLoadUnscoredNoDirectoryReturnsError(t *testing.T) {
	t.Parallel()

	class, err := ParseClassFile(classJSON)
	if err != nil {
		t.Fatalf("ParseClassFile(%q) returned error: %v", classJSON, err)
	}

	err = class.LoadUnscored("testdata/does-not-exist", nil)
	if err == nil {
		t.Fatal("LoadUnscored should return an error for no directory; got nil")
	}
}

func TestLoadGradesNilStudentReturnsError(t *testing.T) {
	t.Parallel()

	class := fakeCalcClass()
	class.studentsByEmail["gstriker@school.edu"] = nil
	dir := t.TempDir()

	gradebookData := `{
    "assignment_date": "20240319",
    "assignment_name": "nil-student",
    "assignment_type": "quiz",
    "assignment_category": "minor",
    "assignment_records": [
        {
            "email": "gstriker@school.edu",
            "grade": 88
        }
    ]
}`

	if err := os.WriteFile(filepath.Join(dir, "quiz-nil-student-20240319.gradebook"), []byte(gradebookData), 0o644); err != nil {
		t.Fatalf("failed writing gradebook fixture: %v", err)
	}

	err := class.LoadGrades(dir, nil)
	if err == nil {
		t.Fatal("LoadGrades should return error for nil student; got nil")
	}
	if !strings.Contains(err.Error(), `student with email "gstriker@school.edu" is nil`) {
		t.Fatalf("unexpected error for nil student: %v", err)
	}
}

func TestLoadGradesNilAssignmentRecordReturnsError(t *testing.T) {
	t.Parallel()

	class := fakeCalcClass()
	dir := t.TempDir()

	gradebookData := `{
    "assignment_date": "20240319",
    "assignment_name": "nil-record",
    "assignment_type": "quiz",
    "assignment_category": "minor",
    "assignment_records": [
        null
    ]
}`

	if err := os.WriteFile(filepath.Join(dir, "quiz-nil-record-20240319.gradebook"), []byte(gradebookData), 0o644); err != nil {
		t.Fatalf("failed writing gradebook fixture: %v", err)
	}

	err := class.LoadGrades(dir, nil)
	if err == nil {
		t.Fatal("LoadGrades should return error for nil assignment record; got nil")
	}
	if !strings.Contains(err.Error(), `nil assignment record at index 0`) {
		t.Fatalf("unexpected error for nil assignment record: %v", err)
	}
}

func TestLoadGradesWithUnicodePrefixAndTermFilter(t *testing.T) {
	t.Parallel()

	class := fakeCalcClass()
	dir := t.TempDir()

	gradebookData := `{
    "assignment_date": "20240319",
    "assignment_name": "unicode-prefix",
    "assignment_type": "quiz",
    "assignment_category": "minor",
    "assignment_records": [
        {
            "email": "gstriker@school.edu",
            "grade": 91
        }
    ]
}`

	fileName := "quiz-na\u00efve-20240319.gradebook"
	if err := os.WriteFile(filepath.Join(dir, fileName), []byte(gradebookData), 0o644); err != nil {
		t.Fatalf("failed writing gradebook fixture: %v", err)
	}

	term := &Term{Start: "20240301", End: "20240331"}
	if err := class.LoadGrades(dir, term); err != nil {
		t.Fatalf("LoadGrades with unicode filename prefix should succeed: %v", err)
	}

	if got := len(class.studentsByEmail["gstriker@school.edu"].gradesByCategory["minor"]); got != 1 {
		t.Fatalf("expected one loaded grade, got %d", got)
	}
}

func writeFinalAndQuizFixtures(t *testing.T) string {
	t.Helper()

	quizData := `{
    "assignment_date": "20210405",
    "assignment_name": "q4-quiz",
    "assignment_type": "quiz",
    "assignment_category": "minor",
    "assignment_records": [
        {
            "email": "gstriker@school.edu",
            "grade": 88
        }
    ]
}`

	finalData := `{
    "assignment_date": "20210608",
    "assignment_name": "year-final",
    "assignment_type": "final",
    "assignment_category": "major",
    "assignment_records": [
        {
            "email": "gstriker@school.edu",
            "grade": 95
        }
    ]
}`

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "quiz-q4-quiz-20210405.gradebook"), []byte(quizData), 0o644); err != nil {
		t.Fatalf("failed writing quiz fixture: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "final-year-final-20210608.gradebook"), []byte(finalData), 0o644); err != nil {
		t.Fatalf("failed writing final fixture: %v", err)
	}

	return dir
}

func TestLoadGradesTermExcludesFinal(t *testing.T) {
	t.Parallel()

	dir := writeFinalAndQuizFixtures(t)
	class := fakeCalcClass()
	class.categoriesByAssignmentType["final"] = "major"

	q4 := &Term{Start: "20210328", End: "20210609"}
	if err := class.LoadGrades(dir, q4); err != nil {
		t.Fatalf("LoadGrades(%q, q4) returned error: %v", dir, err)
	}

	student := class.studentsByEmail["gstriker@school.edu"]
	if got := len(student.gradesByCategory["minor"]); got != 1 {
		t.Fatalf("expected the quiz grade to load, got %d minor grades", got)
	}
	if got := len(student.gradesByCategory["major"]); got != 0 {
		t.Fatalf("expected the final to be excluded when a term is set, got %d major grades", got)
	}
}

func TestLoadGradesNilTermIncludesFinal(t *testing.T) {
	t.Parallel()

	dir := writeFinalAndQuizFixtures(t)
	class := fakeCalcClass()
	class.categoriesByAssignmentType["final"] = "major"

	if err := class.LoadGrades(dir, nil); err != nil {
		t.Fatalf("LoadGrades(%q, nil) returned error: %v", dir, err)
	}

	student := class.studentsByEmail["gstriker@school.edu"]
	if got := len(student.gradesByCategory["minor"]); got != 1 {
		t.Fatalf("expected the quiz grade to load, got %d minor grades", got)
	}
	if got := len(student.gradesByCategory["major"]); got != 1 {
		t.Fatalf("expected the final to be included when term is nil, got %d major grades", got)
	}
}
