//go:generate go-enum  --ptr --marshal --flag --sqlnullstr --sql --names --values --nocomments
package validator

import (
	"github.com/go-playground/validator/v10"
)

// ENUM(EUR, USD)
type Currency string

var CurrencyValidator validator.Func = func(fl validator.FieldLevel) bool {
	// https://github.com/abice/go-enum
	if str, ok := fl.Field().Interface().(string); ok {
		_, err := ParseCurrency(str)
		return err == nil

	}
	return false
}
