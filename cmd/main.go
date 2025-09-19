package main

import (
	"fmt"
	"os"

	"sort_utility/internal/app"
)

func errorExit(err error) {
	_, err = fmt.Fprintln(os.Stderr, err)
	if err != nil {
		panic(err)
	}
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		errorExit(fmt.Errorf("usage: go run main.go <file path>"))
	}
	err := app.RunApp(os.Args...)
	if err != nil {
		errorExit(err)
	}
}
