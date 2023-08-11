package validators

import (
	"time"

	"github.com/go-playground/validator/v10"
)

func BookableDate(lf validator.FieldLevel) bool {
	if date, ok := lf.Field().Interface().(time.Time); ok {
		today := time.Now()
		if today.Year() > date.Year() || today.YearDay() > date.YearDay() {
			return false
		}
	}
	return true
}
