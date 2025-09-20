package args

import (
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectFile  string
		expectError bool
		expectOpts  *KeySort
	}{
		{
			name:       "basic file",
			args:       []string{"test.txt"},
			expectFile: "test.txt",
			expectOpts: &KeySort{SortByColumn: true, ColumnNumber: 1},
		},
		{
			name:       "numeric sort",
			args:       []string{"-n", "test.txt"},
			expectFile: "test.txt",
			expectOpts: &KeySort{SortByColumn: true, ColumnNumber: 1, Numeric: true},
		},
		{
			name:       "reverse sort",
			args:       []string{"-r", "test.txt"},
			expectFile: "test.txt",
			expectOpts: &KeySort{SortByColumn: true, ColumnNumber: 1, Reverse: true},
		},
		{
			name:       "column sort",
			args:       []string{"-k", "3", "test.txt"},
			expectFile: "test.txt",
			expectOpts: &KeySort{SortByColumn: true, ColumnNumber: 3},
		},
		{
			name:        "missing k argument",
			args:        []string{"-k"},
			expectError: true,
		},
		{
			name:        "invalid column number",
			args:        []string{"-k", "0", "test.txt"},
			expectError: true,
		},
		{
			name:        "unknown option",
			args:        []string{"-x", "test.txt"},
			expectError: true,
		},
		{
			name:        "conflicting sort options",
			args:        []string{"-n", "-M", "test.txt"},
			expectError: true,
		},
		{
			name:        "check mode with reverse",
			args:        []string{"-c", "-r", "test.txt"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, opts, err := ParseArgs(tt.args)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if file != tt.expectFile {
				t.Errorf("expected file %q, got %q", tt.expectFile, file)
			}

			if opts.SortByColumn != tt.expectOpts.SortByColumn {
				t.Errorf("expected SortByColumn %v, got %v", tt.expectOpts.SortByColumn, opts.SortByColumn)
			}

			if opts.ColumnNumber != tt.expectOpts.ColumnNumber {
				t.Errorf("expected ColumnNumber %d, got %d", tt.expectOpts.ColumnNumber, opts.ColumnNumber)
			}

			if opts.Numeric != tt.expectOpts.Numeric {
				t.Errorf("expected Numeric %v, got %v", tt.expectOpts.Numeric, opts.Numeric)
			}
		})
	}
}

func TestParseFlag(t *testing.T) {
	tests := []struct {
		name        string
		flags       string
		expectError bool
		checkFunc   func(*KeySort) bool
	}{
		{
			name:      "numeric flag",
			flags:     "n",
			checkFunc: func(ks *KeySort) bool { return ks.Numeric },
		},
		{
			name:      "reverse flag",
			flags:     "r",
			checkFunc: func(ks *KeySort) bool { return ks.Reverse },
		},
		{
			name:      "unique flag",
			flags:     "u",
			checkFunc: func(ks *KeySort) bool { return ks.Unique },
		},
		{
			name:      "month flag",
			flags:     "M",
			checkFunc: func(ks *KeySort) bool { return ks.Month },
		},
		{
			name:      "combined flags",
			flags:     "nr",
			checkFunc: func(ks *KeySort) bool { return ks.Numeric && ks.Reverse },
		},
		{
			name:        "invalid flag",
			flags:       "x",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &KeySort{}
			err := parseFlag(tt.flags, opts)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if !tt.checkFunc(opts) {
				t.Errorf("flag not set correctly")
			}
		})
	}
}
