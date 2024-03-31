package gradebook

import (
	"errors"
	"fmt"
	"gradebook/internal/set"
	"strings"

	"golang.org/x/exp/maps"
)

func zvalErr(zvals []string) error {
	switch len(zvals) {
	case 0:
		return nil
	case 1:
		return fmt.Errorf("gradebook: a field in Class is unset: %s", zvals[0])
	default:
		return fmt.Errorf(
			"gradebooks: fields in Class are unset: %s",
			strings.Join(zvals, ", "),
		)
	}
}

// checkInitialization ensures that a *Class has no dangerous zero values.
func (c *Class) checkInitialization() error {
	zvals := make([]string, 0, 7)
	if c.Name == "" {
		zvals = append(zvals, "Name")
	}
	if c.Terms == nil {
		zvals = append(zvals, "Terms")
	}
	if c.Categories == nil {
		zvals = append(zvals, "Categories")
	}
	if c.PrettyCategories == nil {
		zvals = append(zvals, "PrettyCategories")
	}
	if c.Weights == nil {
		zvals = append(zvals, "Weights")
	}
	if c.Subcategories == nil {
		zvals = append(zvals, "Subcategories")
	}
	if c.Students == nil {
		zvals = append(zvals, "Students")
	}
	return zvalErr(zvals)
}

// checkWeightsSum ensures that c.Weights adds up to 100%.
func (c *Class) checkWeightsSum() error {
	total := 0
	for _, n := range c.Weights {
		total += n
	}
	if total != 100 {
		return errors.New("gradebook: Weights must equal 100%")
	}
	return nil
}

// checkEq returns an error if two sets are not equal or nil if they are.
func checkEq[T comparable](lhs, rhs set.Set[T]) error {
	if !lhs.Equals(rhs) {
		return fmt.Errorf("%s and %s are not equal sets", lhs, rhs)
	}
	return nil
}

// Validate checks whether a *Class is valid. It returns nil if the *Class is
// valid. Otherwise it returns an error containing one more errors from the
// individual checks. Those errors are combined using errors.Join.
func (c *Class) Validate() error {
	cSet := set.New(c.Categories...)
	sSet := set.New(maps.Values(c.Subcategories)...)
	wSet := set.New(maps.Keys(c.Weights)...)
	pcSet := set.New(maps.Keys(c.PrettyCategories)...)
	return errors.Join(
		c.checkInitialization(),
		c.checkWeightsSum(),
		checkEq(cSet, sSet),
		checkEq(cSet, pcSet),
		checkEq(cSet, wSet),
	)
}
