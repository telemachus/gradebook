package gradebook

import (
	"fmt"
	"strings"
)

type set[E comparable] map[E]struct{}

func newSet[E comparable](elems ...E) set[E] {
	s := make(set[E], len(elems))
	for _, el := range elems {
		s[el] = struct{}{}
	}

	return s
}

func (s set[E]) equals(other set[E]) bool {
	if len(s) != len(other) {
		return false
	}

	for el := range s {
		_, ok := other[el]
		if !ok {
			return false
		}
	}

	return true
}

func (s set[E]) String() string {
	elems := make([]string, 0, len(s))
	for el := range s {
		elems = append(elems, fmt.Sprintf("%v", el))
	}

	return fmt.Sprintf("{%s}", strings.Join(elems, ", "))
}
