package app

import (
	"fmt"
	"os"

	p "sort_utility/internal/args"
	f "sort_utility/internal/file"
)

func RunApp(args ...string) error {
	filePath, _, err := p.ParseArgs(args[1:])
	if err != nil {
		return fmt.Errorf("sort:  %s\n", err)
	}
	file, err := f.OpenFile(filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	sortedLines, err := f.SortFile(file)
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
