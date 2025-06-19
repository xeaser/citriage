package dedupe

import (
	"testing"
)

func TestProcess(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name: "Test with both timestamps and consecutive duplicates",
			input: []byte(
				"[2024-05-20T10:00:00Z] Starting build...\n" +
					"[2024-05-20T10:00:01Z] Compiling module A...\n" +
					"[2024-05-20T10:00:02Z] Compiling module A...\n" +
					"[2024-05-20T10:00:03Z] Compiling module A...\n" +
					"[2024-05-20T10:00:04Z] Build complete."),
			expected: "Starting build...\n" +
				"Compiling module A... (repeated 3 times)\n" +
				"Build complete.\n",
		},
		{
			name:     "Test with empty input",
			input:    []byte(""),
			expected: "",
		},
		{
			name: "Test with lines that are empty after timestamp removal",
			input: []byte(
				"[2024-05-20T10:00:00Z] \n" +
					"[2024-05-20T10:00:01Z] \n" +
					"[2024-05-20T10:00:02Z] Some content"),
			expected: "Some content\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := Process(tc.input)

			if actual != tc.expected {
				t.Errorf("Expected output: %s, Actual output: %s", tc.expected, actual)
			}
		})
	}
}
