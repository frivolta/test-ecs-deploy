package api

import (
	db "birdie/db/sqlc"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidDay(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("date_format", validDay)
	// Asserting valid value
	date := "2022-11-06"
	err := validate.Var(date, "date_format")
	assert.Nil(t, err)
	// Asserting error value
	date = "2222-1-3"
	err = validate.Var(date, "date_format")
	assert.Error(t, err)
	// Asserting error char
	date = "2222-10-13l"
	err = validate.Var(date, "date_format")
	assert.Error(t, err)
}

func TestValidPresence(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("presence", validPresence)
	presence := []db.Presence{"MORNING", "AFTERNOON"}
	err := validate.Var(presence, "presence")
	assert.Nil(t, err)
	stra := []string{"INVALID_PRESENCE"}
	err = validate.Var(stra, "presence")
	assert.Error(t, err)
}
