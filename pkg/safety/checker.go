package safety

import (
	"path/filepath"
)

type Checker struct {
	allowedPatterns []string
	deniedPatterns  []string
}

func NewChecker(allowed, denied []string) *Checker {
	return &Checker{
		allowedPatterns: allowed,
		deniedPatterns:  denied,
	}
}

func (sc *Checker) IsAllowed(command string) bool {
	for _, pattern := range sc.deniedPatterns {
		matched, _ := filepath.Match(pattern, command)
		if matched {
			return false
		}
	}
	
	for _, pattern := range sc.allowedPatterns {
		matched, _ := filepath.Match(pattern, command)
		if matched {
			return true
		}
	}
	
	return len(sc.allowedPatterns) == 0
}