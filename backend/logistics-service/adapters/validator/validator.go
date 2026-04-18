package validator

import (
	coreErrors "logistics-service/logistics-service/core/errors"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	v *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{v: validator.New()}
}

func (val *Validator) ValidateStruct(s any) error {
	if err := val.v.Struct(s); err != nil {
		return coreErrors.ErrValidateFailed
	}
	return nil
}