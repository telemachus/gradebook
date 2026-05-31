package gradebook

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestParseClassFileForGradesEqualMock(t *testing.T) {
	t.Parallel()

	want := fakeCalcClass()
	got, err := ParseClassFileForGrades(classJSON)
	if err != nil {
		t.Fatalf("ParseClassFileForGrades(classJSON) = %v; want nil error", err)
	}

	if diff := cmp.Diff(
		want,
		got,
		cmpopts.IgnoreFields(Class{}, "trusted", "domain"),
		cmp.AllowUnexported(Class{}, Student{}),
	); diff != "" {
		t.Errorf("ParseClassFileForGrades(classJSON) mismatch(-want +got):\n%s", diff)
	}
}

func fakeGradesMap() map[string][]float64 {
	gradesByCategories := make(map[string][]float64, 3)
	gradesByCategories["major"] = make([]float64, 0, 50)
	gradesByCategories["minor"] = make([]float64, 0, 50)
	gradesByCategories["cp"] = make([]float64, 0, 50)

	return gradesByCategories
}
