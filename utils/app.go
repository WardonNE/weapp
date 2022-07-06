package utils

import (
	"os"
	"path/filepath"
)

func WorkingPath() (string, error) {
	return os.Getwd()
}

func BasePath() (string, error) {
	executablePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(executablePath), nil
}
