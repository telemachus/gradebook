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
    "assignment_grades": [
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
  the broadest kinds of graded assignment.  Second, the `categories_pretty`
  section maps the items from `categories` to a more readable explanation of
  the category.  Third, the `category_weights` section lists what percentage
  of the final grade each category has.  In addition, there must be
  a `types_to_categories` map of types of graded item to categories.  (Note
  that `categories_pretty` and `category_weights` cannot include any key that
  is not present in `categories`.  In addition, `types_to_categories` cannot
  use any value expect those in `categories`.)  Finally, the `students`
  section is a map of student emails with information about students' first
  and last names.  (Since emails must be unique these work well as a kind of
  primary key.)
+ `class.json` files are a subset of `JSON` with the following format:

```json
{
	"name": "NAME",
	"terms": {
		"TERM1": {
			"start": "DATE-STRING",
			"end": "DATE-STRING"
		},
		"TERM2": {
			"start": "DATE-STRING",
			"end": "DATE-STRING"
		}
	},
	"categories": [ "CATEGORY1", "CATEGORY2", "ETC." ],
	"categories_pretty": {
		"CATEGORY1": "Some human-readable explanation of Category1",
		"CATEGORY2": "Some human-readable explanation of Category2"
	},
	"category_weights": {
		"CATEGORY1": 50,
		"CATEGORY2": 50
	},
	"types_to_categories": {
		"test": "CATEGORY1",
		"essay": "CATEGORY1",
		"quiz": "CATEGORY1",
		"hw": "CATEGORY2",
		"cp": "CATEGORY2"
	},
	"students": {
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
