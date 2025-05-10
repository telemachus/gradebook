package gradebook

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// UnmarshalCalcClass unmarshals a class.json file into a pointer to Class.
// Unlike UnmarshalClass, this function creates the Grades map needed to store
// grades.
func UnmarshalCalcClass(classFile string) (*Class, error) {
	data, err := os.ReadFile(filepath.Clean(classFile))
	if err != nil {
		return nil, err
	}

	var class Class
	err = json.Unmarshal(data, &class)
	if err != nil {
		return nil, err
	}

	for name, student := range class.Students {
		gradesByCategories := make(map[string][]*float64, len(class.Categories))
		for _, cat := range class.Categories {
			// TODO: how should I decide the capacity here?
			gradesByCategories[cat] = make([]*float64, 0, 25)
		}

		student.Grades = gradesByCategories
		class.Students[name] = student
	}

	return &class, nil
}
