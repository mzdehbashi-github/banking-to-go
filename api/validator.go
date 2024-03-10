package api

import (
	"gopsql/banking/util"
	"log"

	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	log.Println("!!!!!!!!!!!!!!!!!!!!!")
	if currency, ok := fl.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}
