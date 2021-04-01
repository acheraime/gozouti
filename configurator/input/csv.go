package input

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
)

func CSVInputReader(inFile string) ([][]string, error) {
	if _, err := os.Stat(inFile); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("%s does not exist", inFile)
		}
	}

	// Open the file
	file, err := os.Open(inFile)
	defer file.Close()
	if err != nil {
		fmt.Printf("attempt to open csv file %s failed: %v", inFile, err)
		return nil, err
	}

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	rows, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("unable to read from io: %v", err)
		return nil, err
	}

	return rows, nil

}
