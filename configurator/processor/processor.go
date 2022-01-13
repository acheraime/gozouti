package processor

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/acheraime/gozouti/utils"
)

type RedirectResources struct {
	dupes     []map[string]string
	Resources []RedirectResource
}

type RedirectResource struct {
	Name        string
	Regex       string
	Replacement string
	ReWriteHost bool
	URLHost     string
}

type Processor interface {
	DryRun(string) error
	Generate() error
}

func NewRedirectResources(input [][]string, parseURL bool, alias string, URLHost string, hostRewrite bool) (RedirectResources, error) {
	var checkBucket = map[string]string{}
	var resources = new(RedirectResources)
	var res = make([]RedirectResource, 0)
	var dupes = []map[string]string{}

	for i, row := range input {
		rFrom := strings.Trim(row[0], " ")
		rTo := strings.Trim(row[1], " ")

		if rFrom == "" || rFrom == "/" {
			// We cannot redirect from /
			continue
		}

		// Check dupe
		if utils.KeyExists(rFrom, checkBucket) {
			resources.dupes = append(dupes, map[string]string{
				"originResource":      rFrom,
				"destinationResource": rTo,
				"rowNumber":           strconv.Itoa(i + 1),
			})

			continue
		}
		checkBucket[rFrom] = rTo

		// A full url is provided
		// parse it to extract the
		// resource paths
		if parseURL {
			furl, err := utils.ParseURL(rFrom)
			if err != nil {
				fmt.Println("unable to parse " + rFrom)
				continue
			}
			rFrom = furl.Path
			turl, err := utils.ParseURL(rTo)
			if err != nil {
				fmt.Println("unable to parse " + rTo)
				continue
			}
			rTo = turl.Path
		}

		rFrom = utils.AddSlash(rFrom)
		rTo = utils.AddSlash(rTo)
		name := fmt.Sprintf("%s-%s", alias, utils.Sanitize(rFrom))
		res = append(res, RedirectResource{Regex: rFrom, Replacement: rTo, Name: name, URLHost: URLHost, ReWriteHost: hostRewrite})
	}

	resources.Resources = res

	return *resources, nil
}
