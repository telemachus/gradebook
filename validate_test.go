package gradebook_test

import (
	"errors"
	"gradebook"
	"testing"
)

func TestZeroValueFieldsGoodClass(t *testing.T) {
	t.Parallel()

	class := testClass()
	if err := class.ZeroValueFields(); err != nil {
		t.Fatalf("want no error for good class; got %s\n", err)
	}
}

func TestZeroValueFieldsBadClasses(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		transformClass func(c *gradebook.Class)
		want           int
	}{
		"missing Name": {
			transformClass: func(c *gradebook.Class) {
				c.Name = ""
			},
			want: 1,
		},
		"missing Name and Terms": {
			transformClass: func(c *gradebook.Class) {
				c.Name = ""
				c.Terms = nil
			},
			want: 2,
		},
		"missing Name, Terms, and Categories": {
			transformClass: func(c *gradebook.Class) {
				c.Name = ""
				c.Terms = nil
				c.Categories = nil
			},
			want: 3,
		},
		"missing Name, Terms, Categories, and CategoriesPretty": {
			transformClass: func(c *gradebook.Class) {
				c.Name = ""
				c.Terms = nil
				c.Categories = nil
				c.CategoriesPretty = nil
			},
			want: 4,
		},
		"missing Name, Terms, Categories, CategoriesPretty, and CategoryWeights": {
			transformClass: func(c *gradebook.Class) {
				c.Name = ""
				c.Terms = nil
				c.Categories = nil
				c.CategoriesPretty = nil
				c.CategoryWeights = nil
			},
			want: 5,
		},
		"missing Name, Terms, Categories, CategoriesPretty, CategoryWeights, and TypesToCategories": {
			transformClass: func(c *gradebook.Class) {
				c.Name = ""
				c.Terms = nil
				c.Categories = nil
				c.CategoriesPretty = nil
				c.CategoryWeights = nil
				c.TypesToCategories = nil
			},
			want: 6,
		},
		"missing Name, Terms, Categories, CategoriesPretty, CategoryWeights, TypesToCategories, and Students": {
			transformClass: func(c *gradebook.Class) {
				c.Name = ""
				c.Terms = nil
				c.Categories = nil
				c.CategoriesPretty = nil
				c.CategoryWeights = nil
				c.TypesToCategories = nil
				c.Students = nil
			},
			want: 7,
		},
	}

	for msg, tc := range testCases {
		tc := tc

		t.Run(msg, func(t *testing.T) {
			t.Parallel()
			class := testClass()
			if tc.transformClass != nil {
				tc.transformClass(class)
			}
			err := class.ZeroValueFields()

			var zverr *gradebook.ZeroValueError

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if errors.As(err, &zverr) && tc.want != len(zverr.Fields) {
				t.Errorf("expected %d, got %d\n", tc.want, len(zverr.Fields))
			}
		})
	}
}

func TestSumCategoryWeightsGoodClass(t *testing.T) {
	t.Parallel()

	class := testClass()
	if err := class.SumCategoryWeights(); err != nil {
		t.Fatalf("want no error for good class; got %s\n", err)
	}
}

func TestSumCategoryWeightsBadClasses(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		transformClass func(c *gradebook.Class)
	}{
		"no weights": {
			transformClass: func(c *gradebook.Class) {
				c.CategoryWeights = gradebook.CategoryWeights{}
			},
		},
		"weights under 100": {
			transformClass: func(c *gradebook.Class) {
				c.CategoryWeights["major"] = 25
			},
		},
		"weights over 100": {
			transformClass: func(c *gradebook.Class) {
				c.CategoryWeights["major"] = 75
			},
		},
		"weights below 0": {
			transformClass: func(c *gradebook.Class) {
				c.CategoryWeights["major"] = -175
			},
		},
	}

	for msg, tc := range testCases {
		tc := tc

		t.Run(msg, func(t *testing.T) {
			t.Parallel()
			class := testClass()
			if tc.transformClass != nil {
				tc.transformClass(class)
			}

			if err := class.SumCategoryWeights(); err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}

func TestConsistencyGoodClass(t *testing.T) {
	t.Parallel()

	class := testClass()
	if err := class.Consistency(); err != nil {
		t.Fatalf("want no error for good class; got %s\n", err)
	}
}

func TestConsistencyBadClasses(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		transformClass func(c *gradebook.Class)
	}{
		"missing item from CategoriesPretty": {
			transformClass: func(c *gradebook.Class) {
				delete(c.CategoriesPretty, "major")
			},
		},
		"extra item in CategoriesPretty": {
			transformClass: func(c *gradebook.Class) {
				c.CategoriesPretty["random"] = "Random Item"
			},
		},
		"missing item from TypesToCategories": {
			transformClass: func(c *gradebook.Class) {
				delete(c.TypesToCategories, "cp")
			},
		},
		"extra item in TypesToCategories": {
			transformClass: func(c *gradebook.Class) {
				c.TypesToCategories["random"] = "random"
			},
		},
		"missing item from CategoryWeights": {
			transformClass: func(c *gradebook.Class) {
				delete(c.CategoryWeights, "cp")
			},
		},
		"extra item in CategoryWeights": {
			transformClass: func(c *gradebook.Class) {
				c.CategoryWeights["random"] = 0
			},
		},
	}

	for msg, tc := range testCases {
		tc := tc

		t.Run(msg, func(t *testing.T) {
			t.Parallel()
			class := testClass()
			if tc.transformClass != nil {
				tc.transformClass(class)
			}

			if err := class.Consistency(); err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}
