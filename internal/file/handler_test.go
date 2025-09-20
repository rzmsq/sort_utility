package file

import (
	"io"
	"os"
	"sort_utility/internal/args"
	"testing"
)

func TestSortFile(t *testing.T) {
	tests := []struct {
		name            string
		content         string
		options         *args.KeySort
		expectedLines   []string
		expectedOutput  string
		shouldReturnNil bool
	}{
		{
			name:          "basic string sort",
			content:       "cherry\napple\nbanana\n",
			options:       &args.KeySort{SortByColumn: true, ColumnNumber: 1},
			expectedLines: []string{"apple", "banana", "cherry"},
		},
		{
			name:          "numeric sort",
			content:       "10\n2\n1\n",
			options:       &args.KeySort{SortByColumn: true, ColumnNumber: 1, Numeric: true},
			expectedLines: []string{"1", "2", "10"},
		},
		{
			name:          "reverse sort",
			content:       "apple\nbanana\ncherry\n",
			options:       &args.KeySort{SortByColumn: true, ColumnNumber: 1, Reverse: true},
			expectedLines: []string{"cherry", "banana", "apple"},
		},
		{
			name:          "unique sort",
			content:       "apple\napple\nbanana\nbanana\ncherry\n",
			options:       &args.KeySort{SortByColumn: true, ColumnNumber: 1, Unique: true},
			expectedLines: []string{"apple", "banana", "cherry"},
		},
		{
			name:          "month sort",
			content:       "march\njanuary\nfebruary\n",
			options:       &args.KeySort{SortByColumn: true, ColumnNumber: 1, Month: true},
			expectedLines: []string{"january", "february", "march"},
		},
		{
			name:          "human numeric sort",
			content:       "2K\n1M\n500\n",
			options:       &args.KeySort{SortByColumn: true, ColumnNumber: 1, HumanNumeric: true},
			expectedLines: []string{"500", "2K", "1M"},
		},
		{
			name:          "column sort",
			content:       "user1 30 admin\nuser2 25 user\nuser3 35 moderator\n",
			options:       &args.KeySort{SortByColumn: true, ColumnNumber: 2, Numeric: true},
			expectedLines: []string{"user2 25 user", "user1 30 admin", "user3 35 moderator"},
		},
		{
			name:            "check sorted file",
			content:         "apple\nbanana\ncherry\n",
			options:         &args.KeySort{SortByColumn: true, ColumnNumber: 1, IsSorted: true},
			expectedOutput:  "Файл отсортирован\n",
			shouldReturnNil: true,
		},
		{
			name:            "check unsorted file",
			content:         "cherry\napple\nbanana\n",
			options:         &args.KeySort{SortByColumn: true, ColumnNumber: 1, IsSorted: true},
			expectedOutput:  "Файл не отсортирован\n",
			shouldReturnNil: true,
		},
		{
			name:          "empty file",
			content:       "",
			options:       &args.KeySort{SortByColumn: true, ColumnNumber: 1},
			expectedLines: []string{},
		},
		{
			name:          "single line",
			content:       "single\n",
			options:       &args.KeySort{SortByColumn: true, ColumnNumber: 1},
			expectedLines: []string{"single"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpFile, err := os.CreateTemp("", "test_*.txt")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			// Write content
			_, err = tmpFile.WriteString(tt.content)
			if err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}
			tmpFile.Close()

			// Reopen for reading
			file, err := os.Open(tmpFile.Name())
			if err != nil {
				t.Fatalf("Failed to open temp file: %v", err)
			}
			defer file.Close()

			// Capture stdout for IsSorted tests
			if tt.options.IsSorted {
				oldStdout := os.Stdout
				r, w, _ := os.Pipe()
				os.Stdout = w

				result, err := SortFile(file, tt.options)

				w.Close()
				os.Stdout = oldStdout

				if err != nil {
					t.Errorf("SortFile() error = %v", err)
					return
				}

				if result != nil && tt.shouldReturnNil {
					t.Errorf("SortFile() should return nil for IsSorted mode")
				}

				output, _ := io.ReadAll(r)
				if string(output) != tt.expectedOutput {
					t.Errorf("Expected output %q, got %q", tt.expectedOutput, string(output))
				}
				return
			}

			// Regular sorting tests
			result, err := SortFile(file, tt.options)
			if err != nil {
				t.Errorf("SortFile() error = %v", err)
				return
			}

			if len(result) != len(tt.expectedLines) {
				t.Errorf("Expected %d lines, got %d", len(tt.expectedLines), len(result))
				return
			}

			for i, line := range result {
				if line != tt.expectedLines[i] {
					t.Errorf("At index %d: expected %q, got %q", i, tt.expectedLines[i], line)
				}
			}
		})
	}
}

func TestSortByColumn(t *testing.T) {
	tests := []struct {
		name     string
		lines    []string
		options  *args.KeySort
		expected []string
	}{
		{
			name:     "basic string sort",
			lines:    []string{"cherry", "apple", "banana"},
			options:  &args.KeySort{SortByColumn: true, ColumnNumber: 1},
			expected: []string{"apple", "banana", "cherry"},
		},
		{
			name:     "numeric sort",
			lines:    []string{"10", "2", "1"},
			options:  &args.KeySort{SortByColumn: true, ColumnNumber: 1, Numeric: true},
			expected: []string{"1", "2", "10"},
		},
		{
			name:     "month sort",
			lines:    []string{"march", "january", "february"},
			options:  &args.KeySort{SortByColumn: true, ColumnNumber: 1, Month: true},
			expected: []string{"january", "february", "march"},
		},
		{
			name:     "human numeric sort",
			lines:    []string{"2K", "1M", "500"},
			options:  &args.KeySort{SortByColumn: true, ColumnNumber: 1, HumanNumeric: true},
			expected: []string{"500", "2K", "1M"},
		},
		{
			name:     "column sort by second column",
			lines:    []string{"user1 30", "user2 25", "user3 35"},
			options:  &args.KeySort{SortByColumn: true, ColumnNumber: 2, Numeric: true},
			expected: []string{"user2 25", "user1 30", "user3 35"},
		},
		{
			name:     "skip blanks",
			lines:    []string{"  apple", " banana", "cherry"},
			options:  &args.KeySort{SortByColumn: true, ColumnNumber: 1, SkipBlanks: true},
			expected: []string{"  apple", " banana", "cherry"},
		},
		{
			name:     "already sorted - no change",
			lines:    []string{"apple", "banana", "cherry"},
			options:  &args.KeySort{SortByColumn: true, ColumnNumber: 1},
			expected: []string{"apple", "banana", "cherry"},
		},
		{
			name:     "column out of range",
			lines:    []string{"a", "bb ccc", "d"},
			options:  &args.KeySort{SortByColumn: true, ColumnNumber: 3},
			expected: []string{"a", "bb ccc", "d"},
		},
		{
			name:     "empty lines",
			lines:    []string{},
			options:  &args.KeySort{SortByColumn: true, ColumnNumber: 1},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy to avoid modifying test data
			lines := make([]string, len(tt.lines))
			copy(lines, tt.lines)

			sortByColumn(lines, tt.options)

			if len(lines) != len(tt.expected) {
				t.Errorf("Expected %d lines, got %d", len(tt.expected), len(lines))
				return
			}

			for i, line := range lines {
				if line != tt.expected[i] {
					t.Errorf("At index %d: expected %q, got %q", i, tt.expected[i], line)
				}
			}
		})
	}
}

func TestWriteMsg(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "simple message",
			input:    []byte("test message"),
			expected: "test message",
		},
		{
			name:     "message with newline",
			input:    []byte("test\nmessage\n"),
			expected: "test\nmessage\n",
		},
		{
			name:     "empty message",
			input:    []byte(""),
			expected: "",
		},
		{
			name:     "unicode message",
			input:    []byte("тестовое сообщение"),
			expected: "тестовое сообщение",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			writeMsg(tt.input)

			w.Close()
			os.Stdout = oldStdout

			output, _ := io.ReadAll(r)
			if string(output) != tt.expected {
				t.Errorf("Expected output %q, got %q", tt.expected, string(output))
			}
		})
	}
}

// Остальные существующие тесты остаются без изменений...
func TestCompareNumeric(t *testing.T) {
	tests := []struct {
		name     string
		a, b     string
		expected bool
	}{
		{"numeric comparison", "1", "2", true},
		{"float comparison", "1.5", "2.5", true},
		{"reverse numeric", "10", "5", false},
		{"string fallback", "abc", "def", true},
		{"mixed types", "abc", "123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compareNumeric(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("compareNumeric(%q, %q) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestCompareMonth(t *testing.T) {
	tests := []struct {
		name     string
		a, b     string
		expected bool
	}{
		{"short months", "jan", "feb", true},
		{"full months", "january", "february", true},
		{"mixed case", "JAN", "feb", true},
		{"reverse order", "mar", "jan", false},
		{"non-month strings", "abc", "def", true},
		{"month vs non-month", "abc", "jan", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compareMonth(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("compareMonth(%q, %q) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestCompareHumanNumeric(t *testing.T) {
	tests := []struct {
		name     string
		a, b     string
		expected bool
	}{
		{"kilobytes", "1K", "2K", true},
		{"megabytes", "1M", "2M", true},
		{"mixed units", "1K", "1M", true},
		{"gigabytes", "1G", "2G", true},
		{"plain numbers", "100", "200", true},
		{"reverse order", "2K", "1K", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compareHumanNumeric(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("compareHumanNumeric(%q, %q) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestParseHumanNumeric(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{"plain number", "100", 100},
		{"kilobyte", "1K", 1024},
		{"megabyte", "1M", 1024 * 1024},
		{"gigabyte", "1G", 1024 * 1024 * 1024},
		{"decimal", "1.5K", 1.5 * 1024},
		{"empty string", "", 0},
		{"invalid", "abc", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseHumanNumeric(tt.input)
			if result != tt.expected {
				t.Errorf("parseHumanNumeric(%q) = %f, want %f", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "no duplicates",
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "with duplicates",
			input:    []string{"a", "a", "b", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "empty slice",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "single element",
			input:    []string{"a"},
			expected: []string{"a"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeDuplicates(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("expected length %d, got %d", len(tt.expected), len(result))
				return
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("at index %d: expected %q, got %q", i, tt.expected[i], v)
				}
			}
		})
	}
}

func TestReverseSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "multiple elements",
			input:    []string{"a", "b", "c", "d"},
			expected: []string{"d", "c", "b", "a"},
		},
		{
			name:     "two elements",
			input:    []string{"a", "b"},
			expected: []string{"b", "a"},
		},
		{
			name:     "single element",
			input:    []string{"a"},
			expected: []string{"a"},
		},
		{
			name:     "empty slice",
			input:    []string{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy to avoid modifying the test data
			input := make([]string, len(tt.input))
			copy(input, tt.input)

			reverseSlice(input)

			if len(input) != len(tt.expected) {
				t.Errorf("expected length %d, got %d", len(tt.expected), len(input))
				return
			}
			for i, v := range input {
				if v != tt.expected[i] {
					t.Errorf("at index %d: expected %q, got %q", i, tt.expected[i], v)
				}
			}
		})
	}
}

func TestIsSorted(t *testing.T) {
	tests := []struct {
		name     string
		lines    []string
		options  *args.KeySort
		expected bool
	}{
		{
			name:     "sorted strings",
			lines:    []string{"apple", "banana", "cherry"},
			options:  &args.KeySort{SortByColumn: true, ColumnNumber: 1},
			expected: true,
		},
		{
			name:     "unsorted strings",
			lines:    []string{"cherry", "apple", "banana"},
			options:  &args.KeySort{SortByColumn: true, ColumnNumber: 1},
			expected: false,
		},
		{
			name:     "sorted numbers",
			lines:    []string{"1", "2", "3"},
			options:  &args.KeySort{SortByColumn: true, ColumnNumber: 1, Numeric: true},
			expected: true,
		},
		{
			name:     "single line",
			lines:    []string{"single"},
			options:  &args.KeySort{SortByColumn: true, ColumnNumber: 1},
			expected: true,
		},
		{
			name:     "empty lines",
			lines:    []string{},
			options:  &args.KeySort{SortByColumn: true, ColumnNumber: 1},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSorted(tt.lines, tt.options)
			if result != tt.expected {
				t.Errorf("isSorted() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestOpenFile(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Write some content
	_, err = tmpFile.WriteString("test content")
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	tests := []struct {
		name        string
		filepath    string
		expectError bool
	}{
		{
			name:        "existing file",
			filepath:    tmpFile.Name(),
			expectError: false,
		},
		{
			name:        "non-existing file",
			filepath:    "non_existing_file.txt",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := OpenFile(tt.filepath)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				if file != nil {
					file.Close()
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if file == nil {
				t.Errorf("expected file but got nil")
				return
			}

			file.Close()
		})
	}
}
