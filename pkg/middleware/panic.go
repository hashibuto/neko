package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/hashibuto/neko"
)

func PanicRecovery(next neko.Handler) neko.Handler {
	return neko.MakeHandler(func(w http.ResponseWriter, r *http.Request) error {
		defer func() {
			err := recover()
			if err != nil {
				fmt.Printf("Recovered from panic\n%v\n%s\n", err, debug.Stack())
			}
		}()

		return next.ServeHTTP(w, r)
	})
}
