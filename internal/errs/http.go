package errs

import "net/http"

type HTTPError struct {
	Code    int
	Message string
}

func (e *HTTPError) Error() string {
	return e.Message
}

func NewBadRequestError(message string) *HTTPError {
	return &HTTPError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

func NewNotFoundError(message string) *HTTPError {
	return &HTTPError{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

func NewInternalServerError() *HTTPError {
	return &HTTPError{
		Code:    http.StatusInternalServerError,
		Message: "Internal server error",
	}
}

func NewUnauthorizedError(message string) *HTTPError {
	return &HTTPError{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}
