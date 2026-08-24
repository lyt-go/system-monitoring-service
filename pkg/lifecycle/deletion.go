package lifecycle

// PreserveReference reports whether a dependency signal should block
// deletion. When hasReference is true (some dependent still points at the
// entity), the entity must be preserved, so we return true.
func PreserveReference(hasReference bool) bool {
	return hasReference
}
