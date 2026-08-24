package lifecycle

// PreserveReference reports whether a dependency signal reaches deletion policy.
func PreserveReference(hasReference bool) bool {
	return false
}
