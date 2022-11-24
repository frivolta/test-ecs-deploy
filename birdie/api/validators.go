package api

import (
	db "birdie/db/sqlc"
	"birdie/util"
	"github.com/go-playground/validator/v10"
)

var validDay validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if date, ok := fieldLevel.Field().Interface().(string); ok {
		_, e := util.ConvertDate(date)
		if e != nil {
			return false
		}
	}
	return true
}

var validPresence validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if _, ok := fieldLevel.Field().Interface().([]db.Presence); ok {
		return true
	}
	return false
}
