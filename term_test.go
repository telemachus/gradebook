package gradebook_test

import (
	"testing"

	"github.com/telemachus/gradebook"
)

func TestTermIncludesValid(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		term    *gradebook.Term
		dateStr string
		want    bool
	}{
		"term should include start date": {
			dateStr: "20230907",
			term: &gradebook.Term{
				Start: "20230907",
				End:   "20231011",
			},
			want: true,
		},
		"term should include end date": {
			dateStr: "20231011",
			term: &gradebook.Term{
				Start: "20230907",
				End:   "20231011",
			},
			want: true,
		},
		"term should include a date between start and end": {
			dateStr: "20230923",
			term: &gradebook.Term{
				Start: "20230907",
				End:   "20231011",
			},
			want: true,
		},
		"term should exclude a date before start": {
			dateStr: "20220907",
			term: &gradebook.Term{
				Start: "20230907",
				End:   "20231011",
			},
			want: false,
		},
		"term should exclude a date after end": {
			dateStr: "20231207",
			term: &gradebook.Term{
				Start: "20230907",
				End:   "20231011",
			},
			want: false,
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			got := tc.term.Includes(tc.dateStr)
			if got != tc.want {
				t.Errorf("%#v.Includes(%s) returns %t; want %t", tc.term, tc.dateStr, got, tc.want)
			}
		})
	}
}

func testTerm(start, end string) *gradebook.Term {
	return &gradebook.Term{
		Start: start,
		End:   end,
	}
}

func loadGradesForTerm(t *testing.T, term *gradebook.Term) *gradebook.Student {
	t.Helper()

	class, err := gradebook.UnmarshalCalcClass("testdata/class.json")
	if err != nil {
		t.Fatalf("failed to unmarshal class: %v", err)
	}

	if err := class.LoadGrades("testdata/term", term); err != nil {
		t.Fatalf("LoadGrades failed: %v", err)
	}

	return class.StudentsByEmail["gstriker@school.edu"]
}

var loadGradesWithTermFilterCases = map[string]struct {
	term             *gradebook.Term
	expectMinorGrade bool
}{
	"loads grades within term": {
		term:             testTerm("20240301", "20240331"),
		expectMinorGrade: true,
	},
	"excludes grades outside term": {
		term:             testTerm("20250101", "20250131"),
		expectMinorGrade: false,
	},
	"nil term loads all grades": {
		term:             nil,
		expectMinorGrade: true,
	},
}

func TestLoadGradesWithTermFilter(t *testing.T) {
	t.Parallel()

	for name, tt := range loadGradesWithTermFilterCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			student := loadGradesForTerm(t, tt.term)
			hasMinorGrade := len(student.GradesByCategory["minor"]) > 0
			if hasMinorGrade != tt.expectMinorGrade {
				t.Errorf("minor-grade presence = %t; want %t", hasMinorGrade, tt.expectMinorGrade)
			}
		})
	}
}

var loadGradesBoundaryCases = map[string]struct {
	termStart        string
	termEnd          string
	expectedCategory string
	shouldLoad       bool
}{
	"march file on start date should be included": {
		termStart:        "20240315",
		termEnd:          "20240331",
		expectedCategory: "minor",
		shouldLoad:       true,
	},
	"march file on end date should be included": {
		termStart:        "20240301",
		termEnd:          "20240315",
		expectedCategory: "minor",
		shouldLoad:       true,
	},
	"march file before start should be excluded": {
		termStart:        "20240320",
		termEnd:          "20240331",
		expectedCategory: "minor",
		shouldLoad:       false,
	},
	"march file after end should be excluded": {
		termStart:        "20240301",
		termEnd:          "20240314",
		expectedCategory: "minor",
		shouldLoad:       false,
	},
	"april file on start date should be included": {
		termStart:        "20240415",
		termEnd:          "20240430",
		expectedCategory: "major",
		shouldLoad:       true,
	},
	"april file on end date should be included": {
		termStart:        "20240401",
		termEnd:          "20240415",
		expectedCategory: "major",
		shouldLoad:       true,
	},
	"april file before start should be excluded": {
		termStart:        "20240420",
		termEnd:          "20240430",
		expectedCategory: "major",
		shouldLoad:       false,
	},
	"april file after end should be excluded": {
		termStart:        "20240401",
		termEnd:          "20240414",
		expectedCategory: "major",
		shouldLoad:       false,
	},
}

func TestLoadGradesTermBoundaries(t *testing.T) {
	t.Parallel()

	for name, tt := range loadGradesBoundaryCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			student := loadGradesForTerm(t, testTerm(tt.termStart, tt.termEnd))
			hasGrade := len(student.GradesByCategory[tt.expectedCategory]) > 0

			if tt.shouldLoad && !hasGrade {
				t.Errorf("expected grades in %s for term %s to %s", tt.expectedCategory, tt.termStart, tt.termEnd)
			}
			if !tt.shouldLoad && hasGrade {
				t.Errorf("expected no grades in %s for term %s to %s", tt.expectedCategory, tt.termStart, tt.termEnd)
			}
		})
	}
}
