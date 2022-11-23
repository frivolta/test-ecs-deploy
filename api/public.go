package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type publicResponse struct {
	Ok string `json:"ok"`
}

func (server *Server) publicRoute(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]string{"status": "200", "title": "Health OK", "detail": time.Now().String()})
}
