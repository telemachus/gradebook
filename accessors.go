package gradebook

import (
	"maps"
	"slices"
)

func (c *Class) termsForLookup() map[string]*Term {
	if domain := c.trustedDomain(); domain != nil {
		return domain.termsByID
	}

	return c.termsByID
}

func (c *Class) categoriesByTypeForLookup() map[string]string {
	if domain := c.trustedDomain(); domain != nil {
		return domain.categoriesByType
	}

	return c.categoriesByAssignmentType
}

func (c *Class) labelsByCategoryForLookup() map[string]string {
	if domain := c.trustedDomain(); domain != nil {
		return domain.labelsByCategory
	}

	return c.labelsByAssignmentCategory
}

func (c *Class) weightsByCategoryForLookup() map[string]int {
	if domain := c.trustedDomain(); domain != nil {
		return domain.weightsByCategory
	}

	return c.weightsByAssignmentCategory
}

// TermByID returns the term for a given term ID.
func (c *Class) TermByID(id string) (*Term, bool) {
	if c == nil {
		return nil, false
	}

	terms := c.termsForLookup()
	if terms == nil {
		return nil, false
	}

	term, ok := terms[id]
	if !ok || term == nil {
		return nil, false
	}

	if c.trustedDomain() != nil {
		return &Term{Start: term.Start, End: term.End}, true
	}

	return term, true
}

// AssignmentCategoryForType returns the assignment category for a given type.
func (c *Class) AssignmentCategoryForType(assignmentType string) (string, bool) {
	if c == nil {
		return "", false
	}

	categoriesByType := c.categoriesByTypeForLookup()
	if categoriesByType == nil {
		return "", false
	}

	category, ok := categoriesByType[assignmentType]
	if !ok {
		return "", false
	}

	return category, true
}

// AssignmentLabelByCategory returns the label for a given assignment category.
func (c *Class) AssignmentLabelByCategory(category string) (string, bool) {
	if c == nil {
		return "", false
	}

	labelsByCategory := c.labelsByCategoryForLookup()
	if labelsByCategory == nil {
		return "", false
	}

	label, ok := labelsByCategory[category]
	if !ok {
		return "", false
	}

	return label, true
}

// AssignmentTypes returns all known assignment types.
func (c *Class) AssignmentTypes() []string {
	if c == nil {
		return []string{}
	}

	categoriesByType := c.categoriesByTypeForLookup()
	if len(categoriesByType) == 0 {
		return []string{}
	}

	return slices.Collect(maps.Keys(categoriesByType))
}

// AssignmentCategoryWeights returns a copy of category weights.
func (c *Class) AssignmentCategoryWeights() WeightsByAssignmentCategory {
	if c == nil {
		return WeightsByAssignmentCategory{}
	}

	weightsByCategory := c.weightsByCategoryForLookup()
	if len(weightsByCategory) == 0 {
		return WeightsByAssignmentCategory{}
	}

	return maps.Clone(weightsByCategory)
}
