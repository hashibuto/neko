package middleware

import (
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
		if err != nil {
			fmt.Printf("Error occurred during request\n%v\n", err)
		}

		statusCode := neko.GetStatusCode(w, err)
		dur := time.Now().UTC().Sub(now).String()
		fmt.Printf("%s - %s %s %d (%s)\n", now.Format(time.RFC3339Nano), r.Method, r.URL.Path, statusCode, dur)

		return err
	})
}
