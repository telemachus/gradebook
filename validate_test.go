package gradebook_test

import (
	"gradebook"
	"testing"
)

func TestValid(t *testing.T) {
	t.Parallel()

	class := fakeClass()
	if err := class.Validate(); err != nil {
		t.Errorf("class.ZeroValues() = %v; want no error", err)
	}
}

func TestInitializationInvalid(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		transformClass func(c *gradebook.Class)
	}{
		"class.Name unset": {
			transformClass: func(c *gradebook.Class) {
				c.Name = ""
			},
		},
		"class.Terms unset": {
			transformClass: func(c *gradebook.Class) {
				c.Terms = nil
			},
		},
		"class.Categories unset": {
			transformClass: func(c *gradebook.Class) {
				c.Categories = nil
			},
		},
		"class.PrettyCategories unset": {
			transformClass: func(c *gradebook.Class) {
				c.PrettyCategories = nil
			},
		},
		"class.Weights unset": {
			transformClass: func(c *gradebook.Class) {
				c.Weights = nil
			},
		},
		"class.Subcategories unset": {
			transformClass: func(c *gradebook.Class) {
				c.Subcategories = nil
			},
		},
		"class.Students unset": {
			transformClass: func(c *gradebook.Class) {
				c.Students = nil
			},
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			class := fakeClass()
			if tc.transformClass != nil {
				tc.transformClass(class)
			}

			if err := class.Validate(); err == nil {
				t.Error("c.Validate() returns nil; want error for zero value(s)")
			}
		})
	}
}

func TestWeightsSumInvalid(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		transformClass func(c *gradebook.Class)
	}{
		"no Weights": {
			transformClass: func(c *gradebook.Class) {
				c.Weights = gradebook.Weights{}
			},
		},
		"Weights under 100": {
			transformClass: func(c *gradebook.Class) {
				c.Weights["major"] = 25
			},
		},
		"Weights over 100": {
			transformClass: func(c *gradebook.Class) {
				c.Weights["major"] = 75
			},
		},
		"Weights below 0": {
			transformClass: func(c *gradebook.Class) {
				c.Weights["major"] = -175
			},
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			class := fakeClass()
			if tc.transformClass != nil {
				tc.transformClass(class)
			}

			if err := class.Validate(); err == nil {
				t.Error("class.Validate() returns nil; want error for incorrect Weights")
			}
		})
	}
}

func TestSetEqualityInvalid(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		transformClass func(c *gradebook.Class)
	}{
		"missing items from Categories": {
			transformClass: func(c *gradebook.Class) {
				clear(c.Categories)
			},
		},
		"extra item in Categories": {
			transformClass: func(c *gradebook.Class) {
				c.Categories = append(c.Categories, "random")
			},
		},
		"missing item from Subcategories": {
			transformClass: func(c *gradebook.Class) {
				delete(c.Subcategories, "cp")
			},
		},
		"extra item in Subcategories": {
			transformClass: func(c *gradebook.Class) {
				c.Subcategories["random"] = "random"
			},
		},
		"missing item from PrettyCategories": {
			transformClass: func(c *gradebook.Class) {
				delete(c.PrettyCategories, "major")
			},
		},
		"extra item in PrettyCategories": {
			transformClass: func(c *gradebook.Class) {
				c.PrettyCategories["random"] = "Random Item"
			},
		},
		"missing item from Weights": {
			transformClass: func(c *gradebook.Class) {
				delete(c.Weights, "cp")
			},
		},
		"extra item in Weights": {
			transformClass: func(c *gradebook.Class) {
				c.Weights["random"] = 0
			},
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			class := fakeClass()
			if tc.transformClass != nil {
				tc.transformClass(class)
			}

			if err := class.Validate(); err == nil {
				t.Errorf("class.Validate() returns nil; want error for %s", msg)
			}
		})
	}
}
