package array

import "strings"

func In(a []string, b string) bool {
	for _, t := range a {
		if strings.EqualFold(t, b) {
			return true
		}
	}

	return false
}
