package args

import (
	"fmt"
	"strconv"
)

// KeySort Sort key
type KeySort struct {
	columnNumber int  // Num for Sort by column
	sortByColumn bool // Sort by column
	numeric      bool // Sort by numeric value (strings are interpreted as numbers).
	reverse      bool // Reverse the sense of comparison.
	unique       bool // Do not output duplicate strings (only unique ones)
	month        bool // Flag for comparison by month name
	skipBlanks   bool // Skip leading blanks when finding end
	isSorted     bool // Check if the data is sorted
	humanNumeric bool // Flag for sorting by human-readable
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
					return "", nil, fmt.Errorf("option requires an argument -- k")
				}

				columnNum, err := strconv.Atoi(args[i+1])
				if err != nil {
					return "", nil, fmt.Errorf("invalid number: %s", args[i+1])
				}
				options.sortByColumn = true
				options.columnNumber = columnNum
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
			return "", nil, fmt.Errorf("cannot read: %s: No such file or directory", filePath)
		}
	}
	return filePath, options, nil
}

// parseFlag A simple function for setting flags in a KeySort structure
func parseFlag(keys string, optionSort *KeySort) error {
	for _, key := range keys {
		switch key {
		case 'k':
			return fmt.Errorf("option -k requires an argument")
		case 'n':
			optionSort.numeric = true
		case 'r':
			optionSort.reverse = true
		case 'u':
			optionSort.unique = true
		case 'M':
			optionSort.month = true
		case 'b':
			optionSort.skipBlanks = true
		case 'c':
			optionSort.isSorted = true
		case 'h':
			optionSort.humanNumeric = true
		default:
			return fmt.Errorf("unknown option: %c", key)
		}
	}
	return nil
}
