# TODO

+ Write what's needed for a `gb-calc` command. See also `TODO-CALC.md` for
  other thoughts.
    + Must provide functions to read `*.gradebook` files and add their content
      to `Student` objects.
    + `Student` must have a field to store grades by categories, probably
      a map of strings to slices of pointers to `float64` values.
    + Must provide functions to average `float64` values and logic to get
      an all-in average for a student given a set of `Class.Weights`.
+ Separate out material from gradebook.go into other files.
