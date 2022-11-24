package api

import (
	db "birdie/db/sqlc"
	"birdie/errorsmanager"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

type kidResponse struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

func newKidResponse(k *db.Kid) kidResponse {
	return kidResponse{
		ID:      k.ID,
		Name:    k.Name,
		Surname: k.Surname,
	}
}

type createKidRequest struct {
	Name    string `json:"name" binding:"required,min=1,max=50,alpha"`
	Surname string `json:"surname" binding:"required,min=1,max=50,alpha"`
}

func (server *Server) createKidRoute(ctx *gin.Context) {
	var req createKidRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorsmanager.ErrorResponse(err))
		return
	}
	arg := db.CreateKidParams{
		Name:    req.Name,
		Surname: req.Surname,
	}

	k, e := server.store.CreateKid(ctx, arg)
	if e != nil {
		ctx.JSON(http.StatusInternalServerError, errorsmanager.ErrorResponse(e))
	}
	ctx.JSON(http.StatusOK, newKidResponse(&k))
}

func (server *Server) getAllKidsRoute(ctx *gin.Context) {
	ks, e := server.store.GetAllKids(ctx)
	if e != nil {
		ctx.JSON(http.StatusInternalServerError, errorsmanager.ErrorResponse(e))
		return
	}
	ctx.JSON(http.StatusOK, ks)
}

type getKidRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// Get a single kid
func (server *Server) getKidRoute(ctx *gin.Context) {
	var req getKidRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorsmanager.ErrorResponse(err))
		return
	}
	k, e := server.store.GetKid(ctx, req.ID)
	if e != nil {
		if e == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorsmanager.ErrorResponse(e))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorsmanager.ErrorResponse(e))
		return
	}

	ctx.JSON(http.StatusOK, newKidResponse(&k))
}
