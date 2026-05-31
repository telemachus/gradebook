package gradebook

import (
	"cmp"
	"maps"
	"slices"
)

func (c *Class) assignmentCategoriesForSorting() []string {
	if c == nil {
		return nil
	}

	if domain := c.trustedDomain(); domain != nil {
		return domain.assignmentCategories
	}

	return c.assignmentCategories
}

func (c *Class) labelsByCategoryForSorting() map[string]string {
	if c == nil {
		return nil
	}

	if domain := c.trustedDomain(); domain != nil {
		return domain.labelsByCategory
	}

	return c.labelsByAssignmentCategory
}

func (c *Class) studentsByEmailForSorting() map[string]*Student {
	if c == nil {
		return nil
	}

	if domain := c.trustedDomain(); domain != nil {
		return domain.studentsByEmail
	}

	return c.studentsByEmail
}

// AssignmentCategoriesSortedByLabel returns a slice of categories sorted by label.
func (c *Class) AssignmentCategoriesSortedByLabel() []string {
	categories := c.assignmentCategoriesForSorting()
	if len(categories) == 0 {
		return []string{}
	}

	labelsByCategory := c.labelsByCategoryForSorting()
	categories = slices.Clone(categories)
	slices.SortFunc(categories, func(catA, catB string) int {
		labelA := labelsByCategory[catA]
		labelB := labelsByCategory[catB]

		return cmp.Compare(labelA, labelB)
	})

	return categories
}

// EmailsSortedByStudentName returns a slice of student emails sorted by student name.
func (c *Class) EmailsSortedByStudentName() []string {
	studentsByEmail := c.studentsByEmailForSorting()
	if len(studentsByEmail) == 0 {
		return []string{}
	}

	emails := slices.Collect(maps.Keys(studentsByEmail))
	slices.SortFunc(emails, func(emailA, emailB string) int {
		studentA := studentsByEmail[emailA]
		studentB := studentsByEmail[emailB]

		return cmpStudent(studentA, studentB)
	})

	return emails
}

// StudentsSortedByName returns a slice of students sorted by last and first name.
func (c *Class) StudentsSortedByName() []*Student {
	studentsByEmail := c.studentsByEmailForSorting()
	students := make([]*Student, 0, len(studentsByEmail))
	for _, student := range studentsByEmail {
		students = append(students, student)
	}
	slices.SortFunc(students, cmpStudent)

	if c.trustedDomain() != nil {
		cloned := make([]*Student, 0, len(students))
		for _, student := range students {
			cloned = append(cloned, cloneStudent(student))
		}

		return cloned
	}

	return students
}

func cmpStudent(studentA, studentB *Student) int {
	if studentA.lastName == studentB.lastName {
		return cmp.Compare(studentA.firstName, studentB.firstName)
	}

	return cmp.Compare(studentA.lastName, studentB.lastName)
}
