package file

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	p "sort_utility/internal/args"
)

func writeMsg(str []byte) {
	_, err := os.Stdout.Write(str)
	if err != nil {
		panic(err)
	}
}

// OpenFile attempts to open the file at the given filepath
// Returns the opened \*os.File and a nil error on success
// If the file does not exist, returns a wrapped ErrFileNotFound error
// For other errors, returns a wrapped error with context
func OpenFile(filepath string) (*os.File, error) {
	file, err := os.Open(filepath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("cannot read: %s: %w", filepath, p.ErrFileNotFound)
		}
		return nil, fmt.Errorf("failed to open file %s: %w", filepath, err)
	}
	return file, nil
}

// SortFile reads lines from the provided file and sorts them based on the given options
// It supports checking if the file is already sorted, sorting by column, removing duplicates,
// and reversing the result. Returns the sorted lines or an error
func SortFile(file *os.File, options *p.KeySort) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if options.IsSorted {
		if isSorted(lines, options) {
			writeMsg([]byte("Файл отсортирован\n"))

		} else {
			writeMsg([]byte("Файл не отсортирован\n"))
		}
		return nil, nil
	}

	sortByColumn(lines, options)

	if options.Unique {
		lines = removeDuplicates(lines)
	}

	if options.Reverse {
		reverseSlice(lines)
	}

	return lines, nil
}

func isSorted(lines []string, options *p.KeySort) bool {
	if len(lines) <= 1 {
		return true
	}

	for i := 0; i < len(lines)-1; i++ {
		fieldsI := strings.Fields(lines[i])
		fieldsJ := strings.Fields(lines[i+1])

		colIndex := options.ColumnNumber - 1
		if colIndex >= len(fieldsI) || colIndex >= len(fieldsJ) {
			if lines[i] > lines[i+1] {
				return false
			}
			continue
		}

		valueI := fieldsI[colIndex]
		valueJ := fieldsJ[colIndex]

		if options.SkipBlanks {
			valueI = strings.TrimSpace(valueI)
			valueJ = strings.TrimSpace(valueJ)
		}

		var greater bool
		if options.Numeric {
			greater = !compareNumeric(valueI, valueJ) && valueI != valueJ
		} else if options.Month {
			greater = !compareMonth(valueI, valueJ) && valueI != valueJ
		} else if options.HumanNumeric {
			greater = !compareHumanNumeric(valueI, valueJ) && valueI != valueJ
		} else {
			greater = valueI > valueJ
		}

		if greater {
			return false
		}
	}
	return true
}

func sortByColumn(lines []string, options *p.KeySort) {
	if isSorted(lines, options) {
		return
	}

	sort.Slice(lines, func(i, j int) bool {
		fieldsI := strings.Fields(lines[i])
		fieldsJ := strings.Fields(lines[j])

		colIndex := options.ColumnNumber - 1
		if colIndex >= len(fieldsI) || colIndex >= len(fieldsJ) {
			return lines[i] < lines[j]
		}

		valueI := fieldsI[colIndex]
		valueJ := fieldsJ[colIndex]

		if options.SkipBlanks {
			valueI = strings.TrimSpace(valueI)
			valueJ = strings.TrimSpace(valueJ)
		}

		if options.Numeric {
			return compareNumeric(valueI, valueJ)
		}

		if options.Month {
			return compareMonth(valueI, valueJ)
		}

		if options.HumanNumeric {
			return compareHumanNumeric(valueI, valueJ)
		}

		return valueI < valueJ
	})
}

func compareNumeric(a, b string) bool {
	numA, errA := strconv.ParseFloat(a, 64)
	numB, errB := strconv.ParseFloat(b, 64)

	if errA != nil && errB != nil {
		return a < b
	}
	if errA != nil {
		return false
	}
	if errB != nil {
		return true
	}

	return numA < numB
}

func compareMonth(a, b string) bool {
	monthOrder := map[string]int{
		"jan": 1, "january": 1,
		"feb": 2, "february": 2,
		"mar": 3, "march": 3,
		"apr": 4, "april": 4,
		"may": 5,
		"jun": 6, "june": 6,
		"jul": 7, "july": 7,
		"aug": 8, "august": 8,
		"sep": 9, "september": 9,
		"oct": 10, "october": 10,
		"nov": 11, "november": 11,
		"dec": 12, "december": 12,
	}

	orderA, okA := monthOrder[strings.ToLower(a)]
	orderB, okB := monthOrder[strings.ToLower(b)]

	if !okA && !okB {
		return a < b
	}
	if !okA {
		return false
	}
	if !okB {
		return true
	}

	return orderA < orderB
}

func compareHumanNumeric(a, b string) bool {
	valueA := parseHumanNumeric(a)
	valueB := parseHumanNumeric(b)
	return valueA < valueB
}

func parseHumanNumeric(s string) float64 {
	if len(s) == 0 {
		return 0
	}

	suffix := strings.ToUpper(string(s[len(s)-1]))
	var multiplier float64 = 1
	var numStr string

	switch suffix {
	case "K":
		multiplier = 1024
		numStr = s[:len(s)-1]
	case "M":
		multiplier = 1024 * 1024
		numStr = s[:len(s)-1]
	case "G":
		multiplier = 1024 * 1024 * 1024
		numStr = s[:len(s)-1]
	case "T":
		multiplier = 1024 * 1024 * 1024 * 1024
		numStr = s[:len(s)-1]
	default:
		numStr = s
	}

	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0
	}

	return num * multiplier
}

func removeDuplicates(lines []string) []string {
	if len(lines) <= 1 {
		return lines
	}

	result := make([]string, 0, len(lines))
	result = append(result, lines[0])

	for i := 1; i < len(lines); i++ {
		if lines[i] != lines[i-1] {
			result = append(result, lines[i])
		}
	}

	return result
}

func reverseSlice(lines []string) {
	for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
		lines[i], lines[j] = lines[j], lines[i]
	}
}
