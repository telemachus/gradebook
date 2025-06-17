# TODO for gb-calc

Here's an outline of how `gradebook calc` works in its Python version.

1. Parse arguments.
1. Load `class.json` file.
1. Create a dict of student objects.
1. Load grades into the student objects.  If the user passes a term filter,
   then only grades for that term are loaded.  Otherwise, all grades are
   loaded.
1. Display grades.

How should this look in Go?

1. Student objects become student structs.  These structs should have the
   following fields: `FirstName`, `LastName`, `Email`, `Grades`.  The first
   three should be strings; `Grades` should be a map of strings to slices of
   pointers to float. The strings will represent categories of grades taken
   from the `class.json` file.  The slices will hold grades for each category.
1. I should create a `New` function, just for `gb-calc`, that creates the
   `Grades` map for each `Student` struct in the gradebook.  This function
   should wrap the `Unmarshal` function that other gradebook-related programs
   will use.  (Other programs won't need a map of `Grades` for students.)
1. Grades in JSON can be `null`.  In Go, I should probably use pointers.
   A `nil` pointer represents a null grade.  If a grade is `null`, then we
   don't add it to the student's grades.
