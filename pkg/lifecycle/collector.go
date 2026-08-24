package lifecycle

import "strings"

func CollectorStopped(status string) bool {
	status = strings.TrimSpace(status)
	return strings.EqualFold(status, "stopped")
}
