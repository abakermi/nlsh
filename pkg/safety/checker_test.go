package safety

import "testing"

func TestSafetyChecker(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		allowed  []string
		denied   []string
		expected bool
	}{
		{
			name:     "allowed command",
			command:  "ls -la",
			allowed:  []string{"ls *"},
			denied:   []string{},
			expected: true,
		},
		{
			name:     "denied command",
			command:  "rm -rf /",
			allowed:  []string{"ls *"},
			denied:   []string{"rm *"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checker := NewChecker(tt.allowed, tt.denied)
			result := checker.IsAllowed(tt.command)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}