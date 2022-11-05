package neko

import (
	"fmt"
	"regexp"
	"strings"
)

var vTokenFinder = regexp.MustCompile("\\{[^{]*\\}")
var vTokenValidator = regexp.MustCompile("^\\{([a-z_]+)(:([a-z]+))?\\}$")

type Route struct{}

func ParseRoute(routePath string, isPrefix bool) (*Route, error) {
	tokens := strings.Split(routePath, "/")
	for _, token := range tokens {
		match := vTokenFinder.MatchString(token)
		if match {
			if token[0] != '{' || token[len(token)-1] != '}' {
				return nil, fmt.Errorf("Route tokens must occupy the entire span between path delimiters")
			}

			matches := vTokenValidator.FindStringSubmatch(token)
			if matches == nil {
				return nil, fmt.Errorf("Route \"%s\" contains a malformed token \"%s\"", routePath, token)
			}

			tokenName := matches[1]
			tokenType := matches[3]
			if tokenType == "" {
				tokenType = "string"
			}
		}
	}

	return &Route{}, nil
}
