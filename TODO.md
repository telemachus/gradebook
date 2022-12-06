# TODO

+ Add validation for the Class struct.  This validation has several points.
    + The `Categories` slice must not be empty.
    + Every value in the `Categories` slice must be a key in the
      `CategoriesPretty` map.
    + Every key in `CategoriesPretty` must be in the `Categories` slice.
    + The values in `CategoriesPretty` must add up to 100.
    + Values in `TypesToCategories` must be present in the `Categories` slice.

  Are those four checks too much for one function?  If so, then the first,
  second, and fourth can be one sub-function, and the fourth can be a second
  sub-function.

  Should I create custom errors?  Probably, but I should also think about this
  some more.
+ The `LoadClass` function should use this validation. That is, the
  `LoadClass` function can fail in multiple ways.  (1) The file may be
  unreadable, (2) the JSON may be invalid, and (3) the Class struct may be
  invalid.  Ideally, the `LoadClass` function (and others) should return
  distinct errors for these situations.  The `LoadClass` function can handle
  (1) and (2) itself, but it should call a validation function to check for
  (3).
+ Add methods that return emails, names, and numbers for a class.
+ Expand the `Student` struct so that it can store grades.
+ Add methods to the `Student` struct so that it can calculate grades.
