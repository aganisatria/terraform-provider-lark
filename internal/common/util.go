package common

import "strings"

// contains checks if string contains any of the substrings.
func contains(s string, substrs ...string) bool {
	for _, substr := range substrs {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}