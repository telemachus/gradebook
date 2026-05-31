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

func parseTrustedClassDomain(c *Class) (*trustedClassDomain, error) {
	if c == nil {
		return nil, errors.New("gradebook: class is nil")
	}

	domain := newTrustedClassDomain(c)

	return domain, errors.Join(
		c.checkInitialization(),
		checkTrustedTerms(domain),
		checkTrustedStudents(domain),
		checkTrustedWeightsSum(domain),
		checkTrustedAssignmentSets(domain),
	)
}

func checkTrustedWeightsSum(domain *trustedClassDomain) error {
	total := 0
	for _, n := range domain.weightsByCategory {
		total += n
	}

	if total != 100 {
		return errors.New(
			"gradebook: weights by assignment category must equal 100",
		)
	}

	return nil
}

func checkTrustedTerms(domain *trustedClassDomain) error {
	errs := make([]error, 0, len(domain.termsByID))
	for id, term := range domain.termsByID {
		if term == nil {
			errs = append(errs, fmt.Errorf("gradebook: term %q is nil", id))

			continue
		}

		start, startErr := time.Parse("20060102", term.Start)
		if startErr != nil {
			errs = append(
				errs,
				fmt.Errorf(
					"gradebook: term %q has invalid start date %q",
					id,
					term.Start,
				),
			)
		}

		end, endErr := time.Parse("20060102", term.End)
		if endErr != nil {
			errs = append(
				errs,
				fmt.Errorf(
					"gradebook: term %q has invalid end date %q",
					id,
					term.End,
				),
			)
		}

		if startErr == nil && endErr == nil && start.After(end) {
			errs = append(
				errs,
				fmt.Errorf(
					"gradebook: term %q start date is after end date",
					id,
				),
			)
		}
	}

	return errors.Join(errs...)
}

func checkTrustedStudents(domain *trustedClassDomain) error {
	errs := make([]error, 0, len(domain.studentsByEmail))
	for email, student := range domain.studentsByEmail {
		if email == "" {
			errs = append(
				errs,
				errors.New("gradebook: student email must not be empty"),
			)
		}
		if email != "" && strings.TrimSpace(email) != email {
			errs = append(
				errs,
				fmt.Errorf(
					"gradebook: student email %q has leading or trailing whitespace",
					email,
				),
			)
		}
		if email != "" && !strings.Contains(email, "@") {
			errs = append(
				errs,
				fmt.Errorf(
					"gradebook: student email %q must contain @",
					email,
				),
			)
		}
		if student == nil {
			errs = append(
				errs,
				fmt.Errorf("gradebook: student %q is nil", email),
			)

			continue
		}
		if strings.TrimSpace(student.firstName) == "" {
			errs = append(
				errs,
				fmt.Errorf(
					"gradebook: student %q first_name must not be empty",
					email,
				),
			)
		}
		if strings.TrimSpace(student.lastName) == "" {
			errs = append(
				errs,
				fmt.Errorf(
					"gradebook: student %q last_name must not be empty",
					email,
				),
			)
		}
	}

	return errors.Join(errs...)
}

func checkTrustedAssignmentSets(domain *trustedClassDomain) error {
	assignmentsSet := set.New(domain.assignmentCategories...)
	categoriesSet := set.New(slices.Collect(maps.Values(domain.categoriesByType))...)
	weightsSet := set.New(slices.Collect(maps.Keys(domain.weightsByCategory))...)
	labelsSet := set.New(slices.Collect(maps.Keys(domain.labelsByCategory))...)

	return errors.Join(
		checkEq(assignmentsSet, categoriesSet),
		checkEq(assignmentsSet, labelsSet),
		checkEq(assignmentsSet, weightsSet),
	)
}
