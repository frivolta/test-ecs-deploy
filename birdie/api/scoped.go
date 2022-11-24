package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type scopedResponse struct {
	Ok string `json:"ok"`
}

func (server *Server) scopedRoute(ctx *gin.Context) {
	var r = scopedResponse{
		Ok: "true",
	}
	ctx.JSON(http.StatusOK, r)
}
