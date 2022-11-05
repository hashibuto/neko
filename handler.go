package neko

import "net/http"

type HandlerFunc func(http.ResponseWriter, *http.Request) error
type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request) error
}

type HandlerStruct struct {
	handlerFunc HandlerFunc
}

func (hs *HandlerStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return hs.handlerFunc(w, r)
}

func MakeHandler(hf HandlerFunc) Handler {
	return &HandlerStruct{
		handlerFunc: hf,
	}
}

func Cascade(finalHandler Handler, middlewares ...Middleware) Handler {
	var curHandler Handler
	for i := range middlewares {
		// Reverse select
		curIndex := len(middlewares) - i - 1
		if i == 0 {
			curHandler = middlewares[curIndex](finalHandler)
		} else {
			curHandler = middlewares[curIndex](curHandler)
		}
	}

	return curHandler
}
