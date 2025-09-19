package app

import (
	"errors"
	"fmt"
	"os"

	p "sort_utility/internal/args"
	f "sort_utility/internal/file"
)

func RunApp(args ...string) error {
	filePath, options, err := p.ParseArgs(args[1:])
	if err != nil {
		return fmt.Errorf("sort: %s", err)
	}
	file, err := f.OpenFile(filePath)
	if err != nil {
		if errors.Is(err, p.ErrFileNotFound) {
			return fmt.Errorf("sort: %w", err)
		}
		return err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	sortedLines, err := f.SortFile(file, options)
	if err != nil {
		return err
	}

	for _, line := range sortedLines {
		_, err = os.Stdout.Write([]byte(line + "\n"))
		if err != nil {
			return err
		}
	}

	return nil
}
