package gradebook

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

var expectedCategories = []string{"major", "minor", "cp"}

func fakeStudentsByEmailSpec() StudentsByEmailSpec {
	students := fakeStudentsByEmail()
	specs := make(StudentsByEmailSpec, len(students))
	for email, student := range students {
		specs[email] = &StudentSpec{
			FirstName: student.firstName,
			LastName:  student.lastName,
		}
	}

	return specs
}

func fakeClassSpec() *ClassSpec {
	return &ClassSpec{
		Name:                        "Lucretius",
		TermsByID:                   fakeTermsByID(),
		AssignmentCategories:        fakeAssignmentCategories(),
		LabelsByAssignmentCategory:  fakeLabelsByAssignmentCategory(),
		WeightsByAssignmentCategory: fakeWeightsByAssignmentCategory(),
		CategoriesByAssignmentType:  fakeCategoriesByAssignmentType(),
		StudentsByEmail:             fakeStudentsByEmailSpec(),
	}
}

func mustBuildClass(t *testing.T, mode InitMode) *Class {
	t.Helper()

	class, err := BuildClass(fakeClassSpec(), mode)
	if err != nil {
		t.Fatalf("BuildClass() returned error: %v", err)
	}

	return class
}

func requireInvalidClassError(t *testing.T, spec *ClassSpec) {
	t.Helper()

	_, err := BuildClass(spec, InitNone)
	if err == nil {
		t.Fatal("BuildClass() should return an error")
	}

	var invalidClassErr *InvalidClassError
	if !errors.As(err, &invalidClassErr) {
		t.Fatalf("BuildClass() error type = %T; want *InvalidClassError", err)
	}
}

func TestBuildClassInitNone(t *testing.T) {
	t.Parallel()

	class := mustBuildClass(t, InitNone)
	student := class.studentsByEmail["gstriker@school.edu"]
	if student == nil {
		t.Fatal("student gstriker@school.edu missing")
	}
	if student.gradesByCategory != nil {
		t.Fatal("InitNone should not initialize GradesByCategory")
	}
	if student.unscoredByCategory != nil {
		t.Fatal("InitNone should not initialize UnscoredByCategory")
	}
}

func TestBuildClassInitGrades(t *testing.T) {
	t.Parallel()

	class := mustBuildClass(t, InitGrades)
	student := class.studentsByEmail["gstriker@school.edu"]
	if student == nil {
		t.Fatal("student gstriker@school.edu missing")
	}
	if student.gradesByCategory == nil {
		t.Fatal("InitGrades should initialize GradesByCategory")
	}
	for _, category := range expectedCategories {
		if _, ok := student.gradesByCategory[category]; !ok {
			t.Fatalf("InitGrades missing category %q", category)
		}
	}
	if student.unscoredByCategory != nil {
		t.Fatal("InitGrades should not initialize UnscoredByCategory")
	}
}

func TestBuildClassInitUnscored(t *testing.T) {
	t.Parallel()

	class := mustBuildClass(t, InitUnscored)
	student := class.studentsByEmail["gstriker@school.edu"]
	if student == nil {
		t.Fatal("student gstriker@school.edu missing")
	}
	if student.unscoredByCategory == nil {
		t.Fatal("InitUnscored should initialize UnscoredByCategory")
	}
	for _, category := range expectedCategories {
		if _, ok := student.unscoredByCategory[category]; !ok {
			t.Fatalf("InitUnscored missing category %q", category)
		}
	}
	if student.gradesByCategory != nil {
		t.Fatal("InitUnscored should not initialize GradesByCategory")
	}
}

func TestBuildClassInvalidClassError(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		requireInvalidClassError(t, nil)
	})

	t.Run("bad-weights", func(t *testing.T) {
		t.Parallel()

		spec := fakeClassSpec()
		spec.WeightsByAssignmentCategory["cp"] = 20
		requireInvalidClassError(t, spec)
	})

	t.Run("bad-term", func(t *testing.T) {
		t.Parallel()

		spec := fakeClassSpec()
		spec.TermsByID["q1"].Start = "2020-09-10"
		requireInvalidClassError(t, spec)
	})

	t.Run("bad-student-email", func(t *testing.T) {
		t.Parallel()

		spec := fakeClassSpec()
		student := spec.StudentsByEmail["gstriker@school.edu"]
		delete(spec.StudentsByEmail, "gstriker@school.edu")
		spec.StudentsByEmail["gstriker.school.edu"] = student
		requireInvalidClassError(t, spec)
	})
}

func TestBuildClassUsesTrustedDomainSnapshot(t *testing.T) {
	t.Parallel()

	class := mustBuildClass(t, InitGrades)
	class.categoriesByAssignmentType["quiz"] = "unknown"
	delete(class.studentsByEmail, "gstriker@school.edu")
	class.assignmentCategories = AssignmentCategories{"major"}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() returned error: %v", err)
	}
	dir := filepath.Join(wd, "testdata", "validgradebooks")

	if err := class.LoadGrades(dir, nil); err != nil {
		t.Fatalf("LoadGrades(%q, nil) returned error: %v", dir, err)
	}
}

func TestBuildClassInvalidInitMode(t *testing.T) {
	t.Parallel()

	_, err := BuildClass(fakeClassSpec(), InitMode(8))
	if err == nil {
		t.Fatal("BuildClass() should return an error for invalid init mode")
	}
}
