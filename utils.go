package neko

import (
	"net/http"
	"reflect"
	"strconv"
)

func ParsePathTokens(r *http.Request) map[string]any {
	mapping := map[string]any{}
	routeVal := r.Context().Value(routeKey)
	route := routeVal.(*Route)
	matches := route.pathRegex.FindStringSubmatch(r.URL.Path)
	if matches != nil {
		for tokenIdx, match := range matches[1:] {
			token := route.pathTokens[tokenIdx]
			if token.kind == reflect.Int {
				// Ignore errors and treat as zero value if incorrect type
				value, _ := strconv.Atoi(match)
				mapping[token.name] = value
			} else {
				mapping[token.name] = match
			}
		}
	}
	return mapping
}
