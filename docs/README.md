# `gradebook`

## Go package to store and calculate grades

This package provides data structures and essential functions necessary to
store and calculate grades for a class.  Grades are stored in `.gradebook`
files, and information about a class is stored in a `class.json` file.  (See
below for more details about both of these filetypes.)

## `.gradebook` files

+ Grades are stored in `.gradebook` files in a directory.  These files have
  names with the following structure: `TYPE-NAME-DATE.gradebook`.  The `TYPE`
  must be found in the class's `class.json` file, but the `NAME` and `DATE`
  portions can be any legal filename.
+ Gradebook files are a subset of `JSON` with the following format:

```json
{
    "assignment_date": "DATE",
    "assignment_name": "NAME",
    "assignment_type": "TYPE",
    "assignment_category": "CATEGORY",
    "assignment_records": [
        {
            "email": "STUDENT_EMAIL",
            "grade": null
        },
        {
            "email": "STUDENT_EMAIL",
            "grade": null
        }
    ]
}
```

## `class.json` files

+ To use the package and tools built on top of this package, a `class.json`
  file must be present in the current working directory.
+ The `class.json` file must define the following essential features of
  a class: the name, terms (i.e., temporal divisions), grading categories,
  grading rubric, and students.  The terms can be overlapping: for example,
  four quarters and two semesters.  The grading categories are distributed in
  four subsections. First, the categories section is a list of brief names of
  the broadest kinds of graded assignment.  Second, the `pretty_categories`
  section maps the items from `categories` to a more readable explanation of
  the category.  (This section was previously labeled `categories_pretty`.)
  Third, the `weights` section lists what percentage of the final grade each
  category has.  (This section was previously labeled `category_weights`.)  In
  addition, there must be a `subcategories` map of types of graded item to
  categories.  (This section was previously labeled `types_to_categories`.)
+ Note that `pretty_categories` and `weights` cannot include any key that is
  not present in `categories`.  In addition, `subcategories` cannot use any
  value expect those in `categories`.  Finally, the `students` section is
  a map of student emails with information about students' first and last
  names. (Since emails must be unique these work well as a kind of primary
  key.)
+ `class.json` files are a subset of `JSON` with the following format:

```json
{
	"name": "NAME",
	"terms_by_id": {
		"TERM1": {
			"start": "DATE-STRING",
			"end": "DATE-STRING"
		},
		"TERM2": {
			"start": "DATE-STRING",
			"end": "DATE-STRING"
		}
	},
	"assignment_categories": [ "CATEGORY1", "CATEGORY2", "ETC." ],
	"labels_by_assignment_category": {
		"CATEGORY1": "Some human-readable version of CATEGORY1",
		"CATEGORY2": "Some human-readable version of CATEGORY2"
	},
	"weights_by_assignment_category": {
		"CATEGORY1": 50,
		"CATEGORY2": 25,
		"CATEGORY3": 25
	},
	"categories_by_assignment_type": {
		"ASSIGNMENT_TYPE1": "CATEGORY1",
		"ASSIGNMENT_TYPE2": "CATEGORY1",
		"ASSIGNMENT_TYPE3": "CATEGORY1",
		"ASSIGNMENT_TYPE4": "CATEGORY2",
		"ASSIGNMENT_TYPE5": "CATEGORY2",
		"ASSIGNMENT_TYPE6": "CATEGORY3"
	},
	"students_by_email": {
		"somestudent@school.edu": {
			"first_name": "SOME",
			"last_name": "STUDENT"
		},
		"anotherstudent@school.edu": {
			"first_name": "ANOTHER",
			"last_name": "STUDENT"
		}
	}
}
```
