package errorsmanager

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "min":
		return "Should be greater than " + fe.Param()
	case "max":
		return "Should be less than " + fe.Param()
	case "date_format":
		return "Should be a valid date YYYY-MM-DD"
	}

	return "Something went wrong validating fields"
}

func MapValidationErrors(ctx *gin.Context, err error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]ErrorMsg, len(ve))
		for i, fe := range ve {
			out[i] = ErrorMsg{
				Field:   fe.Field(),
				Message: getErrorMsg(fe),
			}
		}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": out})
	}
}
