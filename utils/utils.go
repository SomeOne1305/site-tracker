package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/go-playground/validator/v10"
)
func validationErrorsToMap(errs validator.ValidationErrors) map[string]string {
    errors := make(map[string]string)
    for _, e := range errs {
        errors[e.Field()] = e.Error()
    }
    return errors
}

var validate = validator.New()
func ValidateStruct(s interface{}) map[string]string {
    err := validate.Struct(s)
    if err != nil {
        if validationErrs, ok := err.(validator.ValidationErrors); ok {
            return validationErrorsToMap(validationErrs)
        }
        return map[string]string{"error": "Validation failed"}
    }
    return nil
}

func GenerateSixDigitOTP() (string, error) {
	max := big.NewInt(1_000_000) // 0 to 999,999
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP: %w", err)
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}