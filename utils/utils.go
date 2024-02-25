package utils

import (
	"bufio"
	"fmt"
	"os"
)

func ReadFromFile(file string) ([]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %s", file, err)
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file %s: %s", file, err)
	}

	return lines, nil
}

func ReadFromStdin() ([]string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading from stdin: %s", err)
	}

	return lines, nil
}

func IsStringInList(target string, list []string) bool {
	for _, element := range list {
		if element == target {
			return true
		}
	}
	return false
}
