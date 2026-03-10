package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
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
			require.Equal(t, tt.expected, ByteCountSI(tt.b))
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
			require.Equal(t, tt.expected, ByteCountIEC(tt.b))
		})
	}
}
