package neko

import (
	"errors"
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

// GetStatusCode returns the response status code to the present moment in time
func GetStatusCode(w http.ResponseWriter, err error) int {
	rw := w.(*ResponseWriter)
	if rw.WroteHeader() == true {
		return rw.StatusCode()
	}

	var myErr *StatusErr
	if err != nil && errors.As(err, &myErr) {
		return myErr.StatusCode
	}

	if err != nil {
		return 500
	}

	return 200
}

// IsResponseError returns the state of the application response with respect to error at the present time
func IsResponseError(w http.ResponseWriter, err error) bool {
	rw := w.(*ResponseWriter)
	if rw.WroteHeader() == true {
		statusCode := rw.StatusCode()
		if statusCode < 200 || statusCode >= 400 {
			return true
		}

		return false
	}

	var myErr *StatusErr
	if err != nil && errors.As(err, &myErr) {
		if myErr.StatusCode < 200 || myErr.StatusCode >= 400 {
			return true
		}

		return false
	}

	if err != nil {
		return true
	}

	return false
}
