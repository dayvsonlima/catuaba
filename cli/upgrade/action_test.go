package upgrade

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected int
	}{
		{"equal versions", "0.1.5", "0.1.5", 0},
		{"patch greater", "0.1.6", "0.1.5", 1},
		{"patch lesser", "0.1.4", "0.1.5", -1},
		{"minor greater", "0.2.0", "0.1.5", 1},
		{"minor lesser", "0.1.5", "0.2.0", -1},
		{"major greater", "1.0.0", "0.9.9", 1},
		{"major lesser", "0.9.9", "1.0.0", -1},
		{"a has more parts", "0.1.5.1", "0.1.5", 1},
		{"b has more parts", "0.1.5", "0.1.5.1", -1},
		{"extra zeros equal", "1.0.0", "1.0", 0},
		{"both single digit", "2", "1", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compareVersions(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}
