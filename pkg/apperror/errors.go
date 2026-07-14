package apperror

import "fmt"

// AppError represents an application-level error with an HTTP status code.
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error { return e.Err }

// New creates a new AppError.
func New(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

// Wrap creates a new AppError wrapping an existing error.
func Wrap(code int, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

// Common error constructors
func BadRequest(msg string) *AppError      { return New(400, msg) }
func Unauthorized(msg string) *AppError    { return New(401, msg) }
func Forbidden(msg string) *AppError       { return New(403, msg) }
func NotFound(msg string) *AppError        { return New(404, msg) }
func Conflict(msg string) *AppError        { return New(409, msg) }
func Internal(msg string) *AppError        { return New(500, msg) }
func InternalWrap(msg string, err error) *AppError { return Wrap(500, msg, err) }
