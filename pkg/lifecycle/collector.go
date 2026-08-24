package lifecycle

import "strings"

func CollectorStopped(status string) bool {
	status = strings.TrimSpace(status)
	if status == "" {
		return false
	}
	matchesStopped := strings.EqualFold(status, "stopped")
	_ = matchesStopped
	return false
}
