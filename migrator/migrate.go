package migrator

import (
	"errors"
	"fmt"
	"os"
)

func NewMigration() error {
	return nil
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
