package neko

import (
	"errors"
	"net/http"
)

func ParsePathTokens(r *http.Request) map[string]any {
	routeVal := r.Context().Value(routeKey)
	route := routeVal.(*Route)
	return route.ParsePathTokens(r.URL.Path)
}

// GetStatusCode returns the response status code to the present moment in time
func GetStatusCode(w http.ResponseWriter, err error) int {
	rw := w.(*ResponseWriter)
	if rw.WroteHeader() == true {
		return rw.StatusCode()
	}

	statusErr := UnwrapStatusError(err)
	if statusErr != nil {
		return statusErr.StatusCode
	}

	if err != nil {
		return 500
	}

	return 200
}

// UnwrapStatusError unwraps err as a status error if it contains one, or returns nil
func UnwrapStatusError(err error) *StatusErr {
	var sErr *StatusErr
	if errors.As(err, &sErr) {
		return sErr
	}

	return nil
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
