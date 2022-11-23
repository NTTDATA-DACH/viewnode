package utils

import (
	"gotest.tools/v3/assert"
	"testing"
)

func TestByteCountSI(t *testing.T) {
	tests := []struct {
		name     string
		b        int64
		expected string
	}{
		{
			name:     "should convert byte to string according to SI",
			b:        23164899274,
			expected: "23.2 GB",
		},
		{
			name:     "should not convert if value smaller than unit",
			b:        999,
			expected: "999 B",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, ByteCountSI(tt.b), tt.expected)
		})
	}
}

func TestByteCountIEC(t *testing.T) {
	tests := []struct {
		name     string
		b        int64
		expected string
	}{
		{
			name:     "should convert byte to string according to IEC",
			b:        23164899274,
			expected: "21.6 GiB",
		},
		{
			name:     "should not convert if value smaller than unit",
			b:        1023,
			expected: "1023 B",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, ByteCountIEC(tt.b), tt.expected)
		})
	}
}
