package email

import "regexp"

const (
	genderColumn = "pohlaví"
	genderMale   = "muž"
	genderFemale = "žena"
)

var applyGenderMale = regexp.MustCompile(`\(([^|)]*)\|[^|)]*\)`)
var applyGenderFemale = regexp.MustCompile(`\([^|)]*\|([^|)]*)\)`)
