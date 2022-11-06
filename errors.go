package neko

import "fmt"

type StatusErr struct {
	StatusCode int
	ErrString  string
}

func NewStatusErrf(statusCode int, format string, a ...any) *StatusErr {
	return &StatusErr{
		StatusCode: statusCode,
		ErrString:  fmt.Sprintf(format, a...),
	}
}

func (se *StatusErr) Error() string {
	return se.ErrString
}
