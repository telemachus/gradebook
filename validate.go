package gradebook

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"

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

	if c.Name == "" {
		zvals = append(zvals, "Name")
	}
	if c.TermsByID == nil {
		zvals = append(zvals, "TermsByID")
	}
	if c.AssignmentCategories == nil {
		zvals = append(zvals, "AssignmentCategories")
	}
	if c.LabelsByAssignmentCategory == nil {
		zvals = append(zvals, "LabelsByAssignmentCategory")
	}
	if c.WeightsByAssignmentCategory == nil {
		zvals = append(zvals, "WeightsByAssignmentCategory")
	}
	if c.CategoriesByAssignmentType == nil {
		zvals = append(zvals, "CategoriesByAssignmentType")
	}
	if c.StudentsByEmail == nil {
		zvals = append(zvals, "StudentsByEmail")
	}

	return zvalErr(zvals)
}

// checkWeightsSum ensures that c.Weights adds up to 100%.
func (c *Class) checkWeightsSum() error {
	total := 0
	for _, n := range c.WeightsByAssignmentCategory {
		total += n
	}

	if total != 100 {
		return errors.New("gradebook: weights by assignment category must equal 100")
	}

	return nil
}

func (c *Class) checkTerms() error {
	errs := make([]error, 0, len(c.TermsByID))
	for id, term := range c.TermsByID {
		if term == nil {
			errs = append(errs, fmt.Errorf("gradebook: term %q is nil", id))

			continue
		}

		start, startErr := time.Parse("20060102", term.Start)
		if startErr != nil {
			errs = append(errs, fmt.Errorf("gradebook: term %q has invalid start date %q", id, term.Start))
		}

		end, endErr := time.Parse("20060102", term.End)
		if endErr != nil {
			errs = append(errs, fmt.Errorf("gradebook: term %q has invalid end date %q", id, term.End))
		}

		if startErr == nil && endErr == nil && start.After(end) {
			errs = append(errs, fmt.Errorf("gradebook: term %q start date is after end date", id))
		}
	}

	return errors.Join(errs...)
}

func (c *Class) checkStudents() error {
	errs := make([]error, 0, len(c.StudentsByEmail))
	for email, student := range c.StudentsByEmail {
		if email == "" {
			errs = append(errs, errors.New("gradebook: student email must not be empty"))
		}
		if strings.TrimSpace(email) != email {
			errs = append(errs, fmt.Errorf("gradebook: student email %q has leading or trailing whitespace", email))
		}
		if !strings.Contains(email, "@") {
			errs = append(errs, fmt.Errorf("gradebook: student email %q must contain @", email))
		}
		if student == nil {
			errs = append(errs, fmt.Errorf("gradebook: student %q is nil", email))

			continue
		}
		if strings.TrimSpace(student.FirstName) == "" {
			errs = append(errs, fmt.Errorf("gradebook: student %q first_name must not be empty", email))
		}
		if strings.TrimSpace(student.LastName) == "" {
			errs = append(errs, fmt.Errorf("gradebook: student %q last_name must not be empty", email))
		}
	}

	return errors.Join(errs...)
}

// checkEq returns an error if two sets are not equal or nil if they are.
func checkEq[T comparable](lhs, rhs set.Set[T]) error {
	if !lhs.Equals(rhs) {
		return fmt.Errorf("gradebook: %s and %s are not equal sets", lhs, rhs)
	}

	return nil
}

// Validate checks whether a *Class is valid. It returns nil if the *Class is
// valid. Otherwise it returns an error containing one more errors from the
// individual checks. Those errors are combined using errors.Join.
func (c *Class) Validate() error {
	assignmentsSet := set.New(c.AssignmentCategories...)
	categoriesSet := set.New(slices.Collect(maps.Values(c.CategoriesByAssignmentType))...)
	weightsSet := set.New(slices.Collect(maps.Keys(c.WeightsByAssignmentCategory))...)
	labelsSet := set.New(slices.Collect(maps.Keys(c.LabelsByAssignmentCategory))...)

	return errors.Join(
		c.checkInitialization(),
		c.checkTerms(),
		c.checkStudents(),
		c.checkWeightsSum(),
		checkEq(assignmentsSet, categoriesSet),
		checkEq(assignmentsSet, labelsSet),
		checkEq(assignmentsSet, weightsSet),
	)
}
