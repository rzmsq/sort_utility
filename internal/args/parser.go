package args

import "fmt"

func ParseArgs(args []string) (string, string, error) {
	var filePath, options string
	for _, arg := range args {
		if arg[0] == '-' {
			options += arg[1:]
		} else if filePath == "" {
			filePath = arg
		} else {
			return "", "", fmt.Errorf("cannot read: %s: No such file or directory", filePath)
		}
	}
	return filePath, options, nil
}
