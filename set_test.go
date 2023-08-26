package gradebook

import (
	"testing"
)

func TestSetEquals(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		lhs  set[int]
		rhs  set[int]
		want bool
	}{
		"empty sets are equal": {
			lhs:  newSet[int](),
			rhs:  newSet[int](),
			want: true,
		},
		"equal one-item sets are equal": {
			lhs:  newSet(1),
			rhs:  newSet(1),
			want: true,
		},
		"equal multi-item sets are equal": {
			lhs:  newSet(1, 2, 3),
			rhs:  newSet(1, 2, 3),
			want: true,
		},
		"equal sets are equal regardless of declaration duplicates": {
			lhs:  newSet(1, 1, 1, 2, 4),
			rhs:  newSet(2, 2, 2, 4, 4, 1),
			want: true,
		},
		"empty set is unequal to set with elements": {
			lhs:  newSet[int](),
			rhs:  newSet(1, 2),
			want: false,
		},
		"unequal sets are unequal": {
			lhs:  newSet(1, 2, 4),
			rhs:  newSet(1, 2),
			want: false,
		},
	}

	for msg, tt := range tests {
		tt := tt

		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			got := tt.lhs.equals(tt.rhs)
			if got != tt.want {
				t.Errorf("%s.equals(%s) = %t; want %t", tt.lhs, tt.rhs, got, tt.want)
			}
		})
	}
}

func TestSetString(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		set  set[int]
		want string
	}{
		"empty set.String == {}": {
			set:  newSet[int](),
			want: "{}",
		},
		"set(1).String = {1}": {
			set:  newSet(1),
			want: "{1}",
		},
	}

	for msg, tt := range tests {
		tt := tt

		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			got := tt.set.String()
			if got != tt.want {
				t.Errorf("set.String() = %q; want %q", got, tt.want)
			}
		})
	}
}
