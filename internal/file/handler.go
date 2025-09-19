package file

import (
	"bufio"
	"os"
	"sort"

	p "sort_utility/internal/args"
)

func OpenFile(filepath string) (*os.File, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
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

	sort.Strings(lines)
	return lines, nil
}
