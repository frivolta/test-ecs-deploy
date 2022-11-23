package middleware

import (
	db "birdie/db/sqlc"
	"birdie/errorsmanager"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// DBUserMiddleware Check if user is in db otherwise create with default role
func DBUserMiddleware(store *db.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("User").(map[string]interface{})
		uuid := c.MustGet("UUID").(string)
		email := user["email"].(string)

		// Check if user exists in db
		u, e := db.Store.GetUser(*store, c, email)
		if e != nil {
			args := db.CreateUserParams{
				FullName: uuid,
				Role:     db.RoleTEACHER,
				Email:    email,
			}
			u, e = db.Store.CreateUser(*store, c, args)
			if e != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, errorsmanager.ErrorResponse(errors.New("database connection error")))
				return
			}
		}
		c.Set("DBUser", u)
		c.Set("Role", u.Role)
		c.Next()
	}
}
