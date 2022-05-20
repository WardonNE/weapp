package utils

import "os"

func FileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		panic(err)
	}
	return true
}

func IsDir(filepath string) bool {
	fi, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		panic(err)
	}
	return fi.IsDir()
}
