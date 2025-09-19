package file

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	p "sort_utility/internal/args"
)

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

func SortFile(file *os.File, options *p.KeySort) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if options.SortByColumn == true {
		sortByColumn(lines, options)
	} else {
		sort.Strings(lines)
	}

	return lines, nil
}

func sortByColumn(lines []string, options *p.KeySort) {
	sort.Slice(lines, func(i, j int) bool {
		fieldsI := strings.Fields(lines[i])
		fieldsJ := strings.Fields(lines[j])

		colIndex := options.ColumnNumber - 1
		if colIndex >= len(fieldsI) || colIndex >= len(fieldsJ) {
			return lines[i] < lines[j]
		}

		valueI := fieldsI[colIndex]
		valueJ := fieldsJ[colIndex]

		return valueI < valueJ
	})
}
