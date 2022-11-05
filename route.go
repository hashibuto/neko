package neko

import (
	"fmt"
	"regexp"
	"strings"
)

var vTokenFinder = regexp.MustCompile("\\{[^{]*\\}")

func ParseRoute(routePath string, isPrefix bool) (*Route, error) {
	tokens := strings.Split(routePath, "/")
	for _, token := range tokens {
		match := vTokenFinder.Match([]byte(token))
		if match {
			if token[0] != '{' || token[len(token)-1] != '}' {
				return nil, fmt.Errorf("Route tokens must occupy the entire span between path delimiters")
			}
		}
	}
}
