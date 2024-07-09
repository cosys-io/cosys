package cms

import (
	"bufio"
	"errors"
	"os"
)

func getModFile() (string, error) {
	file, err := os.Open("go.mod")
	if err != nil {
		return "", nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		return line[7:], nil
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", errors.New("module not found")
}
