package utils

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

const urlRegexp = `(^https?:\/\/)?(\[[\w:.]+\]|[\w\.-]+)?(\/[\w\.-]*[éàèêëēėîïíīìäöôüÄÖÜ\w\.-]*\/?[\w\._-]*[éàèêëēėîïíīìäöôüÄÖÜ\w\.-]*\/?)$`

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

func UserDesktop() string {
	return filepath.Join(HomeDir(), "Desktop")
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

func SplitURL(uri string) ([]string, error) {
	re := regexp.MustCompile(urlRegexp)
	if !re.Match([]byte(uri)) {
		return nil, fmt.Errorf("unable to parse %s. not a valid url", uri)
	}

	match := re.FindStringSubmatch(uri)
	return match, nil
}

// KeyExists check if an item with key key
// exists in map bucket
func KeyExists(key string, bucket map[string]string) bool {
	if _, dupe := bucket[key]; dupe {
		return true
	}

	return false
}

// ParseURL is a wrapper func around
// net/url.Parse
func ParseURL(uri string) (*url.URL, error) {
	return url.Parse(uri)
}

// AddSlash adds a trailing slash to
// a string path if missing
func AddSlash(path string) string {
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	return path
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

func RemoveAccents(s string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	result, _, _ := transform.String(t, s)

	return result
}

func Sanitize(s string) string {
	re := regexp.MustCompile(`[/_,;\s]+?`)
	o := re.ReplaceAllString(s, "-")

	o = strings.Trim(o, "-")

	return RemoveAccents(o)
}
