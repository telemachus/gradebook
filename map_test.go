package gradebook

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestKeysAndVals(t *testing.T) {
	t.Parallel()

	var nilMap map[string]int

	strSort := cmpopts.SortSlices(func(a, b string) bool {
		return a < b
	})

	intSort := cmpopts.SortSlices(func(a, b int) bool {
		return a < b
	})

	tests := map[string]struct {
		given    map[string]int
		keysWant []string
		valsWant []int
	}{
		"nil map": {
			given:    nilMap,
			keysWant: []string{},
			valsWant: []int{},
		},
		"empty map": {
			given:    map[string]int{},
			keysWant: []string{},
			valsWant: []int{},
		},
		`map[string]int{"a": 1, "b":2, "c": 3}`: {
			given:    map[string]int{"a": 1, "b": 2, "c": 3},
			keysWant: []string{"b", "c", "a"},
			valsWant: []int{2, 3, 1},
		},
	}

	for msg, tt := range tests {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			ks := keys(tt.given)
			vs := vals(tt.given)

			if !cmp.Equal(ks, tt.keysWant, strSort) {
				t.Errorf(
					"keys(%v) = %v; want %v",
					tt.given,
					ks,
					tt.keysWant,
				)
			}

			if !cmp.Equal(vs, tt.valsWant, intSort) {
				t.Errorf(
					"vals(%v) = %v; want %v",
					tt.given,
					vs,
					tt.valsWant,
				)
			}
		})
	}
}
