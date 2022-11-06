package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hashibuto/neko"
)

func RequestLogger(next neko.Handler) neko.Handler {
	return neko.MakeHandler(func(w http.ResponseWriter, r *http.Request) error {
		now := time.Now().UTC()
		fmt.Printf("%s - %s %s\n", now.Format(time.RFC3339Nano), r.Method, r.URL.Path)
		err := next.ServeHTTP(w, r)
		rw := w.(*neko.ResponseWriter)
		var statusCode int
		var se *neko.StatusErr
		if rw.WroteHeader() {
			statusCode = rw.StatusCode()
		} else if err != nil {
			if errors.As(err, &se) {
				statusCode = se.StatusCode
			} else {
				statusCode = 500
			}
		} else {
			statusCode = 200
		}
		dur := time.Now().UTC().Sub(now).String()
		fmt.Printf("%s - %s %s %d (%s)\n", now.Format(time.RFC3339Nano), r.Method, r.URL.Path, statusCode, dur)

		return err
	})
}
