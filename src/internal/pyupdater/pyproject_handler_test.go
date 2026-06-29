package pyupdater

import (
	"testing"
)

func TestWithoutVersion(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "no version specifier",
			input:    []string{"requests"},
			expected: []string{"requests"},
		},
		{
			name:     "with >= version",
			input:    []string{"requests>=2.0.0"},
			expected: []string{"requests"},
		},
		{
			name:     "with == version",
			input:    []string{"requests==2.28.0"},
			expected: []string{"requests"},
		},
		{
			name:     "with != version",
			input:    []string{"requests!=2.28.0"},
			expected: []string{"requests"},
		},
		{
			name:     "with extras and version",
			input:    []string{"py-key-value-aio[valkey]>=1.0.0"},
			expected: []string{"py-key-value-aio[valkey]"},
		},
		{
			name:     "platform marker without version",
			input:    []string{"py-key-value-aio[valkey]; platform_system != 'Windows'"},
			expected: []string{"py-key-value-aio[valkey]; platform_system != 'Windows'"},
		},
		{
			name:     "platform marker with version",
			input:    []string{"py-key-value-aio[valkey]>=1.0.0; platform_system != 'Windows'"},
			expected: []string{"py-key-value-aio[valkey]; platform_system != 'Windows'"},
		},
		{
			name:     "multiple deps with mixed markers",
			input:    []string{"requests>=2.0.0", "flask; python_version >= '3.8'", "boto3==1.26.0"},
			expected: []string{"requests", "flask; python_version >= '3.8'", "boto3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := withoutVersion(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("got %d results, want %d", len(result), len(tt.expected))
				return
			}
			for i, got := range result {
				if got != tt.expected[i] {
					t.Errorf("result[%d] = %q, want %q", i, got, tt.expected[i])
				}
			}
		})
	}
}
