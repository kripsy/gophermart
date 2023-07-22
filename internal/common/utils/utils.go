package utils

// StingContains returns is searchString contain in arrayString.
func StingContains(arrayString []string, searchString string) bool {
	for _, v := range arrayString {
		if v == searchString {
			return true
		}
	}
	return false
}
