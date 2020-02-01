// Package diff lib
package diff

import (
	"fmt"

	r3 "github.com/r3labs/diff"
)

//GetDiffChangelog log
func GetDiffChangelog(oldObj, newObj interface{}) (*r3.Changelog, error) {
	if r3.Changed(oldObj, newObj) {
		changelog, err := r3.Diff(oldObj, newObj)
		if err != nil {
			return nil, err
		}

		return &changelog, nil
	}

	return nil, nil
}

//GetDiffString log
func GetDiffString(oldObj, newObj interface{}) ([]string, error) {
	var results []string

	changelog, err := GetDiffChangelog(oldObj, newObj)
	if err != nil {
		return results, err
	}

	for _, c := range *changelog {
		results = append(results, fmt.Sprintf("%v changed from %v to %v", c.Path, c.From, c.To))
	}

	return results, nil
}
