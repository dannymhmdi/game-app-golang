package richerr

import (
	"fmt"
	"net/http"
)

type Kind int

const (
	KindInvalid Kind = iota + 1
	KindForbidden
	KindNotFound
	KindUnexpected
)

type RichError struct {
	operation  string
	msg        string
	wrappedErr error
	kind       Kind
}

func New() *RichError {
	return &RichError{}
}

func (r RichError) SetOperation(op string) RichError {
	r.operation = op
	return r
}

func (r RichError) SetMsg(msg string) RichError {
	r.msg = msg
	return r
}

func (r RichError) SetWrappedErr(wrappedErr error) RichError {
	r.wrappedErr = wrappedErr
	return r
}

func (r RichError) SetKind(k Kind) RichError {
	r.kind = k
	return r
}

func (r RichError) Error() string {
	return r.msg
}

func CheckTypeErr(err error) (code Kind, msg string, op string) {
	switch err.(type) {
	case RichError:
		fmt.Println("metallica")
		richErr := err.(RichError)
		return richErr.kind, richErr.msg, richErr.operation
	default:
		return 0, err.Error(), "unknown operation"
	}
}

func MapKindToHttpErr(code Kind) int {
	switch code {
	case 1:
		return http.StatusBadRequest
	case 2:
		return http.StatusForbidden
	case 3:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
