package validator

import "testing"

func TestNewValidator(t *testing.T) {
	v := New()
	if v == nil {
		t.Fatal("New() returned nil")
	}
	if v.validate == nil {
		t.Fatal("New() returned a Validator with nil validate field")
	}
}

func TestValidateStruct(t *testing.T) {
	v := New()

	type TestStruct struct {
		Name  string `validate:"required,min=3"`
		Email string `validate:"required,email"`
	}

	// Valid struct should pass.
	valid := &TestStruct{Name: "John", Email: "john@example.com"}
	if err := v.ValidateStruct(valid); err != nil {
		t.Errorf("expected no error for valid struct, got: %v", err)
	}

	// Invalid struct should fail.
	invalid := &TestStruct{Name: "AB", Email: "invalid-email"}
	if err := v.ValidateStruct(invalid); err == nil {
		t.Error("expected error for invalid struct, got nil")
	}
}

func TestValidateStruct_Required(t *testing.T) {
	v := New()

	type TestStruct struct {
		Name string `validate:"required"`
	}

	// Missing required field should fail.
	empty := &TestStruct{Name: ""}
	if err := v.ValidateStruct(empty); err == nil {
		t.Error("expected error for missing required field, got nil")
	}

	// Present required field should pass.
	filled := &TestStruct{Name: "test"}
	if err := v.ValidateStruct(filled); err != nil {
		t.Errorf("expected no error for present required field, got: %v", err)
	}
}

func TestValidateStruct_Email(t *testing.T) {
	v := New()

	type TestStruct struct {
		Email string `validate:"required,email"`
	}

	validEmail := &TestStruct{Email: "user@example.com"}
	if err := v.ValidateStruct(validEmail); err != nil {
		t.Errorf("expected no error for valid email, got: %v", err)
	}

	invalidEmail := &TestStruct{Email: "not-an-email"}
	if err := v.ValidateStruct(invalidEmail); err == nil {
		t.Error("expected error for invalid email, got nil")
	}
}

func TestValidateStruct_Min(t *testing.T) {
	v := New()

	type TestStruct struct {
		Password string `validate:"required,min=8"`
	}

	// Short value should fail min check.
	short := &TestStruct{Password: "short"}
	if err := v.ValidateStruct(short); err == nil {
		t.Error("expected error for value shorter than min, got nil")
	}

	// Exact min length should pass.
	exact := &TestStruct{Password: "longpass"}
	if err := v.ValidateStruct(exact); err != nil {
		t.Errorf("expected no error for value at exact min length, got: %v", err)
	}
}

func TestValidate(t *testing.T) {
	v := New()

	type TestStruct struct {
		Name  string `validate:"required,min=3"`
		Email string `validate:"required,email"`
	}

	// Test that Validate is an alias for ValidateStruct.
	valid := &TestStruct{Name: "Alice", Email: "alice@example.com"}
	if err := v.Validate(valid); err != nil {
		t.Errorf("expected no error for valid struct, got: %v", err)
	}

	invalid := &TestStruct{Name: "Al", Email: "alice@example.com"}
	if err := v.Validate(invalid); err == nil {
		t.Error("expected error for struct failing min constraint, got nil")
	}
}

func TestValidateField(t *testing.T) {
	v := New()

	// Valid field with a valid tag.
	if err := v.ValidateField("hello", "required,min=3"); err != nil {
		t.Errorf("expected no error for valid field, got: %v", err)
	}

	// Invalid field value should fail.
	if err := v.ValidateField("hi", "required,min=3"); err == nil {
		t.Error("expected error for field shorter than min, got nil")
	}

	// Email validation on a field.
	if err := v.ValidateField("user@example.com", "email"); err != nil {
		t.Errorf("expected no error for valid email, got: %v", err)
	}

	if err := v.ValidateField("not-an-email", "email"); err == nil {
		t.Error("expected error for invalid email, got nil")
	}
}
