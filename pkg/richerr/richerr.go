package richerr

import (
	"net/http"
)

type Kind int

const (
	KindInvalid Kind = iota + 1
	KindForbidden
	KindNotFound
	KindUnexpected
	KindUnathorized
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

	richErr, ok := err.(RichError)
	if !ok {
		return KindUnexpected, "internal server error", "unknown"
	}

	if richErr.kind == 0 {
		return CheckTypeErr(richErr.wrappedErr)
	}

	if MapKindToHttpErr(richErr.kind) >= 500 {
		return richErr.kind, "internal server error", richErr.operation
	}
	return richErr.kind, richErr.msg, richErr.operation
}

func MapKindToHttpErr(code Kind) int {
	switch code {
	case KindInvalid:
		return http.StatusUnprocessableEntity
	case KindForbidden:
		return http.StatusForbidden
	case KindNotFound:
		return http.StatusNotFound
	case KindUnexpected:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
