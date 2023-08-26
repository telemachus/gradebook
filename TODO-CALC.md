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
   floats. The strings will represent categories of grades taken from the
   `class.json` file.  The slices will hold grades for each category.
1. Grades in JSON can be `null`.  In Go, I should probably use pointers.
   A `nil` pointer represents a null grade.  If a student has no grade for
   a given assignment, then I will simply leave that item out of that
   student's slice for that category of grade.
