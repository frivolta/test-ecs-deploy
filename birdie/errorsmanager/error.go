package errorsmanager

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ErrorBuilder interface {
	ErrorResponse(ctx *gin.Context)
}

type ConcreteError struct {
	message string
	status  int
	stack   string
	ctx     *gin.Context
}

// print error to context
func (b *ConcreteError) ConcreteResponse() {
	b.ctx.AbortWithStatusJSON(b.status, gin.H{"error": b.message, "stack": b.stack})
	return
}

func getErrorType(err error) *ConcreteError {
	if err == sql.ErrNoRows {
		return NotFoundError(err)
	}
	return InternalServerError(err)
}

func BuildError(err error, ctx *gin.Context) *ConcreteError {
	concreteError := getErrorType(err)
	concreteError.ctx = ctx
	concreteError.ConcreteResponse()
	return concreteError
}

// Builders
func NotFoundError(e error) *ConcreteError {
	err := ConcreteError{
		message: "Resource not found",
		status:  http.StatusNotFound,
		stack:   e.Error(),
	}
	return &err
}
func InternalServerError(e error) *ConcreteError {
	err := ConcreteError{
		message: "Something went wrong.",
		status:  http.StatusInternalServerError,
		stack:   e.Error(),
	}
	return &err
}

func BadRequest(e error) *ConcreteError {
	err := ConcreteError{
		message: "Bad request" + e.Error(),
		status:  0,
		stack:   e.Error(),
	}
	return &err
}

// Legacy manager
func ErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
