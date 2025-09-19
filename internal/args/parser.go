package args

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrMissingArgument = errors.New("option requires an argument")
	ErrInvalidNumber   = errors.New("invalid number")
	ErrUnknownOption   = errors.New("unknown option")
	ErrFileNotFound    = errors.New("no such file or directory")
)

// KeySort Sort key
type KeySort struct {
	ColumnNumber int  // Num for Sort by column
	SortByColumn bool // Sort by column
	Numeric      bool // Sort by numeric value (strings are interpreted as numbers).
	Reverse      bool // Reverse the sense of comparison.
	Unique       bool // Do not output duplicate strings (only unique ones)
	Month        bool // Flag for comparison by month name
	SkipBlanks   bool // Skip leading blanks when finding end
	IsSorted     bool // Check if the data is sorted
	HumanNumeric bool // Flag for sorting by human-readable
}

// ParseArgs Parsing flags and file name
func ParseArgs(args []string) (string, *KeySort, error) {
	var filePath string
	options := &KeySort{}

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if arg[0] == '-' {
			// A separate case for parsing the k flag
			if arg == "-k" {
				if i+1 >= len(args) {
					return "", nil, fmt.Errorf("%w -- k", ErrMissingArgument)
				}

				columnNum, err := strconv.Atoi(args[i+1])
				if err != nil {
					return "", nil, fmt.Errorf("%w: %s", ErrInvalidNumber, args[i+1])
				}
				options.SortByColumn = true
				options.ColumnNumber = columnNum
				i++
			} else {
				err := parseFlag(arg[1:], options)
				if err != nil {
					return "", nil, err
				}
			}
		} else if filePath == "" {
			filePath = arg
		} else {
			return "", nil, fmt.Errorf("cannot read: %s: %w", filePath, ErrFileNotFound)
		}
	}
	return filePath, options, nil
}

// parseFlag A simple function for setting flags in a KeySort structure
func parseFlag(keys string, optionSort *KeySort) error {
	for _, key := range keys {
		switch key {
		case 'k':
			return fmt.Errorf("-k %w", ErrMissingArgument)
		case 'n':
			optionSort.Numeric = true
		case 'r':
			optionSort.Reverse = true
		case 'u':
			optionSort.Unique = true
		case 'M':
			optionSort.Month = true
		case 'b':
			optionSort.SkipBlanks = true
		case 'c':
			optionSort.IsSorted = true
		case 'h':
			optionSort.HumanNumeric = true
		default:
			return fmt.Errorf("%w: %c", ErrUnknownOption, key)
		}
	}
	return nil
}
