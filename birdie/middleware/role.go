package middleware

import (
	db "birdie/db/sqlc"
	"birdie/errorsmanager"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// RoleAuthorizationMiddleware check for the role passed in context, if user has no required role exit and return
func RoleAuthorizationMiddleware(role db.Role) gin.HandlerFunc {
	return func(context *gin.Context) {
		cRole := context.MustGet("Role")
		if cRole != role && cRole != db.RoleADMIN {
			context.AbortWithStatusJSON(http.StatusForbidden, errorsmanager.ErrorResponse(errors.New("you are not authorize to access this resource")))
			return
		}
		context.Next()
	}
}
