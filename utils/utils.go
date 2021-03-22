package utils

import (
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
