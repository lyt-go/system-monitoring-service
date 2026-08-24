package lifecycle

// PreserveThresholdReference reports whether a threshold must be kept
// because it is still referenced by at least one alert. When alerts reference
// the threshold it cannot be deleted without orphaning those alerts.
func PreserveThresholdReference(hasReference bool) bool {
	return hasReference
}
