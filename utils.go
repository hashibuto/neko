package neko

import (
	"errors"
	"net/http"
)

type HandlerWrapper struct {
	originalHandler http.Handler
}

// ServeHTTP calls the original ServeHTTP method while providing a nil error response
func (hw *HandlerWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	hw.originalHandler.ServeHTTP(w, r)
	return nil
}

// ParsePathTokens parses path tokens out of the supplied path, according to the route matching path template
func ParsePathTokens(r *http.Request) map[string]any {
	routeVal := r.Context().Value(routeKey)
	route := routeVal.(*Route)
	return route.ParsePathTokens(r.URL.Path)
}

// GetPathTemplate returns the path template used to route the actual path.  Eg: "/entity/123" would be "/entity/{id}"
func GetPathTemplate(r *http.Request) string {
	routeVal := r.Context().Value(routeKey)
	route := routeVal.(*Route)
	return route.path
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

// WrapStandardHandler wraps a standard http.Handler and returns a neko handler
func WrapStandardHandler(handler http.Handler) *HandlerWrapper {
	return &HandlerWrapper{
		originalHandler: handler,
	}
}
