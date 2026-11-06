package types

import (
	"errors"
	"testing"
)

func TestNewAppError(t *testing.T) {
	err := errors.New("underlying error")
	appErr := NewAppError(ErrCodeNotFound, "resource not found", err)

	if appErr.Code != ErrCodeNotFound {
		t.Errorf("Code = %d; want %d", appErr.Code, ErrCodeNotFound)
	}
	if appErr.Message != "resource not found" {
		t.Errorf("Message = %q; want %q", appErr.Message, "resource not found")
	}
	if !errors.Is(appErr, err) {
		t.Error("AppError should wrap underlying error")
	}
}

func TestNewAppErrorWithoutWrapping(t *testing.T) {
	appErr := NewAppError(ErrCodeBadRequest, "invalid input", nil)
	// Verify no nil pointer dereference
	if appErr == nil {
		t.Error("AppError should not be nil")
	}
	if appErr.Code != ErrCodeBadRequest {
		t.Errorf("Code = %d; want %d", appErr.Code, ErrCodeBadRequest)
	}
}

func TestAppError_Unwrap(t *testing.T) {
	err := errors.New("wrapped error")
	appErr := NewAppError(ErrCodeDatabase, "db error", err)
	unwrapped := appErr.Unwrap()
	if unwrapped != err {
		t.Errorf("Unwrap() = %v; want %v", unwrapped, err)
	}
}
