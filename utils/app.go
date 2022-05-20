package utils

import (
	"os"
	"path/filepath"
)

func WorkingPath() (string, error) {
	executablePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(executablePath), nil
}

func BasePath() (string, error) {
	return os.Getwd()
}
