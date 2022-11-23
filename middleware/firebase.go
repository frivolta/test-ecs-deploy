package middleware

import (
	"birdie/errorsmanager"
	"context"
	"errors"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// AuthMiddleware : to verify all authorized operations
func AuthMiddleware(c *gin.Context) {
	firebaseAuth := c.MustGet("firebaseAuth").(*auth.Client)
	authorizationToken := c.GetHeader("Authorization")
	idToken := strings.TrimSpace(strings.Replace(authorizationToken, "Bearer", "", 1))
	if idToken == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, errorsmanager.ErrorResponse(errors.New("empty token")))
		return
	}
	//verify token
	token, err := firebaseAuth.VerifyIDToken(context.Background(), idToken)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, errorsmanager.ErrorResponse(errors.New("invalid token")))
		return
	}
	c.Set("UUID", token.UID)
	c.Set("User", token.Claims)
	c.Next()
}
