package gradebook

import (
	"fmt"
	"strings"

	"github.com/telemachus/gradebook/internal/set"
)

func zvalErr(zvals []string) error {
	switch len(zvals) {
	case 0:
		return nil
	case 1:
		return fmt.Errorf("gradebook: class field is unset: %s", zvals[0])
	default:
		return fmt.Errorf("gradebook: class fields are unset: %s", strings.Join(zvals, ", "))
	}
}

// checkInitialization ensures that a *Class has no dangerous zero values.
func (c *Class) checkInitialization() error {
	zvals := make([]string, 0, 7)

	if c.name == "" {
		zvals = append(zvals, "Name")
	}
	if c.termsByID == nil {
		zvals = append(zvals, "TermsByID")
	}
	if c.assignmentCategories == nil {
		zvals = append(zvals, "AssignmentCategories")
	}
	if c.labelsByAssignmentCategory == nil {
		zvals = append(zvals, "LabelsByAssignmentCategory")
	}
	if c.weightsByAssignmentCategory == nil {
		zvals = append(zvals, "WeightsByAssignmentCategory")
	}
	if c.categoriesByAssignmentType == nil {
		zvals = append(zvals, "CategoriesByAssignmentType")
	}
	if c.studentsByEmail == nil {
		zvals = append(zvals, "StudentsByEmail")
	}

	return zvalErr(zvals)
}

// checkEq returns an error if two sets are not equal or nil if they are.
func checkEq[T comparable](lhs, rhs set.Set[T]) error {
	if !lhs.Equals(rhs) {
		return fmt.Errorf("gradebook: %s and %s are not equal sets", lhs, rhs)
	}

	return nil
}
