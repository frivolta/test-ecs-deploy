package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type privateResponse struct {
	Ok string `json:"ok"`
}

func (server *Server) privateRoute(ctx *gin.Context) {
	var r = privateResponse{
		Ok: "true",
	}
	ctx.JSON(http.StatusOK, r)
}
