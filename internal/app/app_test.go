package app

import (
	"os"
	"testing"
)

func TestRunApp(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write test content
	content := "cherry\napple\nbanana\n"
	_, err = tmpFile.WriteString(content)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "basic sort",
			args:        []string{"program", tmpFile.Name()},
			expectError: false,
		},
		{
			name:        "numeric sort",
			args:        []string{"program", "-n", tmpFile.Name()},
			expectError: false,
		},
		{
			name:        "missing file",
			args:        []string{"program", "non_existing.txt"},
			expectError: true,
		},
		{
			name:        "invalid arguments",
			args:        []string{"program", "-x", tmpFile.Name()},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RunApp(tt.args...)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
