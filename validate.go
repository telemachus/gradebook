package gradebook

import (
	"errors"
	"fmt"
	"strings"
)

// ErrCategoryWeights signals that a class does not have a valid grading rubric.
//
// The total of a CategoryWeights in a Class must equal 100.
var ErrCategoryWeights = errors.New("the total of CategoryWeights must equal 100")

// ZeroValueError stores a list of fields in a class that still have their
// zero values.
type ZeroValueError struct {
	Fields []string
}

func (e *ZeroValueError) Error() string {
	switch len(e.Fields) {
	case 0:
		panic("no fields have their zero value")
	case 1:
		return fmt.Sprintf("%q has its zero value", e.Fields[0])
	default:
		return fmt.Sprintf("%q have their zero value", strings.Join(e.Fields, ", "))
	}
}

// ConsistencyError signals that a Class has incoherent definitions of types of
// assignment.
type ConsistencyError struct {
	Conflicts []string
}

func (e *ConsistencyError) Error() string {
	switch len(e.Conflicts) {
	case 0:
		panic("no conflicts")
	default:
		return strings.Join(e.Conflicts, ": ")
	}
}

// ZeroValueFields signals that a Class contains unset fields.
func (c *Class) ZeroValueFields() error {
	zeroValueFields := []string{}
	if c.Name == "" {
		zeroValueFields = append(zeroValueFields, "Name")
	}

	if len(c.Terms) < 1 {
		zeroValueFields = append(zeroValueFields, "Terms")
	}

	if len(c.Categories) < 1 {
		zeroValueFields = append(zeroValueFields, "Categories")
	}

	if len(c.CategoriesPretty) < 1 {
		zeroValueFields = append(zeroValueFields, "CategoriesPretty")
	}

	if len(c.CategoryWeights) < 1 {
		zeroValueFields = append(zeroValueFields, "CategoryWeights")
	}

	if len(c.TypesToCategories) < 1 {
		zeroValueFields = append(zeroValueFields, "TypesToCategories")
	}

	if len(c.Students) < 1 {
		zeroValueFields = append(zeroValueFields, "Students")
	}

	if len(zeroValueFields) > 0 {
		return &ZeroValueError{
			Fields: zeroValueFields,
		}
	}

	return nil
}

// SumCategoryWeights returns ErrCategoryWeights when CategoryWeights does not
// add up to 100.
func (c *Class) SumCategoryWeights() error {
	total := 0

	for _, n := range c.CategoryWeights {
		total += n
	}

	if total != 100 {
		return ErrCategoryWeights
	}

	return nil
}

// Consistency returns ConsistencyError when a Class has incoherent fields.
func (c *Class) Consistency() error {
	conflicts := []string{}
	categories := make(map[string]struct{}, len(c.Categories))
	typesToCategoriesValues := make(map[string]struct{}, len(c.TypesToCategories))

	for _, category := range c.Categories {
		categories[category] = struct{}{}
	}

	for _, value := range c.TypesToCategories {
		typesToCategoriesValues[value] = struct{}{}
	}

	// Every item in Categories must be a key in CategoriesPretty.
	for _, category := range c.Categories {
		if _, ok := c.CategoriesPretty[category]; !ok {
			conflict := fmt.Sprintf("%q is in Categories but not in CategoriesPretty", category)
			conflicts = append(conflicts, conflict)
		}
	}

	// Every key in CategoriesPretty must be in Categories.
	for category := range c.CategoriesPretty {
		if _, ok := categories[category]; !ok {
			conflict := fmt.Sprintf("%q is in CategoriesPretty but not in Categories", category)
			conflicts = append(conflicts, conflict)
		}
	}

	// Every value in TypesToCategories must be in Categories.
	for _, category := range c.TypesToCategories {
		if _, ok := categories[category]; !ok {
			conflict := fmt.Sprintf("%q is in TypesToCategories but not in Categories", category)
			conflicts = append(conflicts, conflict)
		}
	}

	// Every item in Categories must be a value in TypesToCategories.
	for _, category := range c.Categories {
		if _, ok := typesToCategoriesValues[category]; !ok {
			conflict := fmt.Sprintf("%q is in Categories but not in TypesToCategories", category)
			conflicts = append(conflicts, conflict)
		}
	}

	// Every key in CategoryWeights must be in Categories.
	for category := range c.CategoryWeights {
		if _, ok := categories[category]; !ok {
			conflict := fmt.Sprintf("%q is in CategoryWeights but not in Categories", category)
			conflicts = append(conflicts, conflict)
		}
	}

	// Every item in Categories must be a key in CategoryWeights.
	for _, category := range c.Categories {
		if _, ok := c.CategoryWeights[category]; !ok {
			conflict := fmt.Sprintf("%q is in Categories but not in CategoryWeights", category)
			conflicts = append(conflicts, conflict)
		}
	}

	if len(conflicts) > 0 {
		return &ConsistencyError{
			Conflicts: conflicts,
		}
	}

	return nil
}
