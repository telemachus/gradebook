package gradebook_test

import (
	"gradebook"
	"testing"
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
