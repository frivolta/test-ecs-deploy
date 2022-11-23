package api

import (
	utils "birdie/api/utils"
	db "birdie/db/sqlc"
	"birdie/errorsmanager"
	"birdie/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func createCarnetPayload(server *Server, ctx *gin.Context) (utils.KidCarnetInfo, error) {
	// Get all required info from the db
	// Get carnet response and transform it into []db.Carnet
	var allCarnets []db.Carnet
	gcr, e := server.store.GetAllCarnets(ctx)
	if e != nil {
		return nil, e
	}
	for _, v := range gcr {
		allCarnets = append(allCarnets, db.Carnet{
			ID:       v.ID,
			Date:     v.Date,
			Quantity: v.Quantity,
			KidID:    v.KidID.Int64,
		})
	}

	// Get all the kids in the db
	allKids, e := server.store.GetAllKids(ctx)
	if e != nil {
		return nil, e
	}

	// Get all the kid notes in the db
	allKidNotes, e := server.store.GetAllKidNotes(ctx)
	if e != nil {
		return nil, e
	}

	// Init the CUtil to generate the meaningful payload
	cu := utils.NewCUtil(allCarnets)
	cu.ToKidCarnet()
	cu.ToKidCarnetInfo(allKidNotes, allKids)
	return cu.Kci, nil
}

type createCarnetRequest struct {
	Date     string `json:"date" binding:"required,date_format"`
	Quantity int32  `uri:"quantity" binding:"required,min=1"`
}
type createCarnetParams struct {
	KidId int64 `uri:"kid_id" binding:"required,min=1"`
}

func (server *Server) createCarnetRoute(ctx *gin.Context) {
	var par createCarnetParams
	var req createCarnetRequest
	// Bind uri
	if err := ctx.ShouldBindUri(&par); err != nil {
		errorsmanager.MapValidationErrors(ctx, err)
		return
	}
	// Bind request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorsmanager.MapValidationErrors(ctx, err)
		return
	}

	// Check if kid exists
	_, e := server.store.GetKid(ctx, par.KidId)
	if e != nil {
		ne := fmt.Errorf("kid with id %d does not exists", par.KidId)
		errorsmanager.BuildError(ne, ctx)
		return
	}

	d, _ := util.ConvertDate(req.Date)

	arg := db.CreateCarnetParams{
		Date:     d,
		Quantity: req.Quantity,
		KidID:    par.KidId,
	}
	c, e := server.store.CreateCarnet(ctx, arg)
	if e != nil {
		errorsmanager.BuildError(e, ctx)
		return
	}
	ctx.JSON(http.StatusOK, c)
}

type updateCarnetRequest struct {
	Date     string `json:"date" binding:"required,date_format"`
	Quantity int32  `uri:"quantity" binding:"required,min=1"`
}
type updateCarnetParam struct {
	CarnetId int64 `uri:"carnet_id" binding:"required,min=1"`
}

func (server *Server) updateCarnetRoute(ctx *gin.Context) {
	var req updateCarnetRequest
	var param updateCarnetParam
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorsmanager.MapValidationErrors(ctx, err)
		return
	}
	if err := ctx.ShouldBindUri(&param); err != nil {
		errorsmanager.MapValidationErrors(ctx, err)
		return
	}
	_, e := server.store.GetCarnet(ctx, param.CarnetId)
	if e != nil {
		errorsmanager.BuildError(e, ctx)
		return
	}
	d, _ := util.ConvertDate(req.Date)

	// User can only update note text
	arg := db.UpdateCarnetParams{
		ID:       param.CarnetId,
		Date:     d,
		Quantity: req.Quantity,
	}
	uc, e := server.store.UpdateCarnet(ctx, arg)
	if e != nil {
		errorsmanager.BuildError(e, ctx)
		return
	}
	ctx.JSON(http.StatusOK, uc)
}

func (server *Server) getAllCarnetsRoute(ctx *gin.Context) {
	payload, e := createCarnetPayload(server, ctx)
	if e != nil {
		errorsmanager.BuildError(e, ctx)
		return
	}
	ctx.JSON(http.StatusOK, payload)
}
