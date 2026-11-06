package validator

import "github.com/go-playground/validator/v10"

// Validator wraps validator.Validate to provide a framework-level API.
type Validator struct {
	validate *validator.Validate
}

// New creates a new Validator instance.
func New() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

// ValidateStruct validates all fields of the given struct.
func (v *Validator) ValidateStruct(s interface{}) error {
	return v.validate.Struct(s)
}

// Validate is an alias for ValidateStruct.
func (v *Validator) Validate(i interface{}) error {
	return v.validate.Struct(i)
}

// ValidateField validates a single field value against a tag string.
func (v *Validator) ValidateField(field interface{}, tag string) error {
	return v.validate.Var(field, tag)
}
