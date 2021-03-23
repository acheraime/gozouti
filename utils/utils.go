package utils

import (
	"errors"
	"fmt"
	"os"
)

func CheckDir(dir string) error {
	finfo, err := os.Stat(dir)
	if err != nil {
		return err
	}

	if !finfo.IsDir() {
		return fmt.Errorf("%s is not a directory", dir)
	}

	return nil
}

func HomeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}

	return os.Getenv("USERPROFILE")

}

func ReadFile(inFile string) (*os.File, error) {
	if _, err := os.Stat(inFile); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("%s does not exist", inFile)
		}
	}
	// Open the file
	file, err := os.Open(inFile)
	if err != nil {
		fmt.Printf("attempt to open csv file %s failed: %v", inFile, err)
		return nil, err
	}

	return file, nil
}
