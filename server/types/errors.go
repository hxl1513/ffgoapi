package types

const (
	ErrCodeBadRequest     = 4001
	ErrCodeUnauthorized   = 4002
	ErrCodeForbidden      = 4003
	ErrCodeNotFound       = 4004
	ErrCodeValidation     = 4005
	ErrCodeInternalServer = 5001
	ErrCodeDatabase       = 5002
	ErrCodeCache          = 5003
)

// AppError is a structured application error that carries an HTTP-style code,
// a human-readable message, and an optional wrapped cause error.
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

// NewAppError creates a new AppError. When err is non-nil it is stored
// directly so that errors.Is and pointer-equality checks both succeed.
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Error satisfies the error interface.
func (e *AppError) Error() string {
	return e.Message
}

// Unwrap returns the underlying wrapped error so that errors.Is / errors.As
// traverse the chain.
func (e *AppError) Unwrap() error {
	return e.Err
}
