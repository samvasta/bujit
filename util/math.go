package util

// Absolute value for int
func AbsI(an int) int {
	if an < 0 {
		return -an
	}
	return an
}

// Absolute value for int64
func AbsI64(an int64) int64 {
	if an < 0 {
		return -an
	}
	return an
}
