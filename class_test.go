package gradebook

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func testStudent(firstName, lastName string) *Student {
	return &Student{
		firstName: firstName,
		lastName:  lastName,
	}
}

var studentsSortedByNameCases = map[string]struct {
	students map[string]*Student
	want     []*Student
}{
	"empty class returns empty slice": {
		students: map[string]*Student{},
		want:     []*Student{},
	},
	"single student returns single student": {
		students: map[string]*Student{
			"alice@example.com": testStudent("Alice", "Anderson"),
		},
		want: []*Student{
			testStudent("Alice", "Anderson"),
		},
	},
	"multiple students sorted by last name": {
		students: map[string]*Student{
			"charlie@example.com": testStudent("Charlie", "Chen"),
			"alice@example.com":   testStudent("Alice", "Anderson"),
			"bob@example.com":     testStudent("Bob", "Baker"),
		},
		want: []*Student{
			testStudent("Alice", "Anderson"),
			testStudent("Bob", "Baker"),
			testStudent("Charlie", "Chen"),
		},
	},
	"same last name sorted by first name": {
		students: map[string]*Student{
			"bob.smith@example.com":     testStudent("Bob", "Smith"),
			"alice.smith@example.com":   testStudent("Alice", "Smith"),
			"charlie.smith@example.com": testStudent("Charlie", "Smith"),
		},
		want: []*Student{
			testStudent("Alice", "Smith"),
			testStudent("Bob", "Smith"),
			testStudent("Charlie", "Smith"),
		},
	},
	"mixed sorting - last name priority, then first name": {
		students: map[string]*Student{
			"bob.young@example.com":      testStudent("Bob", "Young"),
			"alice.young@example.com":    testStudent("Alice", "Young"),
			"charlie.smith@example.com":  testStudent("Charlie", "Smith"),
			"david.anderson@example.com": testStudent("David", "Anderson"),
		},
		want: []*Student{
			testStudent("David", "Anderson"),
			testStudent("Charlie", "Smith"),
			testStudent("Alice", "Young"),
			testStudent("Bob", "Young"),
		},
	},
}

func TestStudentsSortedByName(t *testing.T) {
	t.Parallel()

	for name, tt := range studentsSortedByNameCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			class := &Class{
				studentsByEmail: tt.students,
			}

			got := class.StudentsSortedByName()

			if diff := cmp.Diff(tt.want, got, cmp.AllowUnexported(Student{})); diff != "" {
				t.Errorf("StudentsSortedByName() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

var emailsSortedByNameCases = map[string]struct {
	students map[string]*Student
	want     []string
}{
	"empty class returns empty slice": {
		students: map[string]*Student{},
		want:     []string{},
	},
	"single student returns single email": {
		students: map[string]*Student{
			"alice@example.com": testStudent("Alice", "Anderson"),
		},
		want: []string{"alice@example.com"},
	},
	"multiple students sorted by last name": {
		students: map[string]*Student{
			"charlie@example.com": testStudent("Charlie", "Chen"),
			"alice@example.com":   testStudent("Alice", "Anderson"),
			"bob@example.com":     testStudent("Bob", "Baker"),
		},
		want: []string{
			"alice@example.com",
			"bob@example.com",
			"charlie@example.com",
		},
	},
	"same last name sorted by first name": {
		students: map[string]*Student{
			"bob.smith@example.com":     testStudent("Bob", "Smith"),
			"alice.smith@example.com":   testStudent("Alice", "Smith"),
			"charlie.smith@example.com": testStudent("Charlie", "Smith"),
		},
		want: []string{
			"alice.smith@example.com",
			"bob.smith@example.com",
			"charlie.smith@example.com",
		},
	},
	"mixed sorting - last name priority, then first name": {
		students: map[string]*Student{
			"bob.young@example.com":      testStudent("Bob", "Young"),
			"alice.young@example.com":    testStudent("Alice", "Young"),
			"charlie.smith@example.com":  testStudent("Charlie", "Smith"),
			"david.anderson@example.com": testStudent("David", "Anderson"),
		},
		want: []string{
			"david.anderson@example.com",
			"charlie.smith@example.com",
			"alice.young@example.com",
			"bob.young@example.com",
		},
	},
}

func TestEmailsSortedByStudentName(t *testing.T) {
	t.Parallel()

	for name, tt := range emailsSortedByNameCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			class := &Class{
				studentsByEmail: tt.students,
			}

			got := class.EmailsSortedByStudentName()

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("EmailsSortedByStudentName() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

//nolint:funlen // Keeping tests together matters more than function length.
func TestAssignmentCategoriesSortedByLabel(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		categories AssignmentCategories
		labels     LabelsByAssignmentCategory
		want       []string
	}{
		"empty categories returns empty slice": {
			categories: AssignmentCategories{},
			labels:     LabelsByAssignmentCategory{},
			want:       []string{},
		},
		"single category returns single category": {
			categories: AssignmentCategories{"major"},
			labels: LabelsByAssignmentCategory{
				"major": "Major Assessments",
			},
			want: []string{"major"},
		},
		"multiple categories sorted by label": {
			categories: AssignmentCategories{"major", "cp", "minor"},
			labels: LabelsByAssignmentCategory{
				"major": "Major Assessments",
				"minor": "Daily Work",
				"cp":    "Class Participation",
			},
			want: []string{"cp", "minor", "major"},
		},
		"categories with same label prefix sorted alphabetically": {
			categories: AssignmentCategories{"quiz", "exam", "project"},
			labels: LabelsByAssignmentCategory{
				"quiz":    "Assessment: Quiz",
				"exam":    "Assessment: Exam",
				"project": "Assessment: Project",
			},
			want: []string{"exam", "project", "quiz"},
		},
		"mixed label sorting": {
			categories: AssignmentCategories{"final", "hw", "participation", "midterm"},
			labels: LabelsByAssignmentCategory{
				"final":         "Final Exam",
				"hw":            "Homework",
				"participation": "Class Participation",
				"midterm":       "Midterm Exam",
			},
			want: []string{"participation", "final", "hw", "midterm"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			class := &Class{
				assignmentCategories:       tt.categories,
				labelsByAssignmentCategory: tt.labels,
			}

			got := class.AssignmentCategoriesSortedByLabel()

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("AssignmentCategoriesSortedByLabel() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
