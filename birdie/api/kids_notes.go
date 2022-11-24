package api

import (
	db "birdie/db/sqlc"
	"birdie/errorsmanager"
	"birdie/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type createKidNoteRequest struct {
	Note     string        `json:"note" binding:"min=0,max=250"`
	Date     string        `json:"date" binding:"required,date_format"`
	Presence []db.Presence `json:"presence" binding:"required,presence"`
	HasMeal  *bool         `json:"has_meal" binding:"required"`
}

type createKidNoteParams struct {
	KidId int64 `uri:"kid_id" binding:"required,min=1"`
}

func (server *Server) createKidNoteRoute(ctx *gin.Context) {
	var par createKidNoteParams
	var req createKidNoteRequest
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

	//Check if kid exists
	k, e := server.store.GetKid(ctx, par.KidId)
	if e != nil {
		errorsmanager.BuildError(e, ctx)
		return
	}
	d, _ := util.ConvertDate(req.Date)

	// Check if kid note is already present for that day
	ns, err := server.store.GetKidNotesByPeriod(ctx, db.GetKidNotesByPeriodParams{
		Date:   d,
		Date_2: d,
	})
	for _, v := range ns {
		if v.KidID.Int64 == par.KidId {
			newE := fmt.Errorf("kid already has a note for this day %d", par.KidId)
			errorsmanager.BuildError(newE, ctx)
			return
		}
	}

	// Create note
	arg := db.CreateKidNoteParams{
		Note:     req.Note,
		KidID:    k.ID,
		Presence: req.Presence,
		HasMeal:  *req.HasMeal,
		Date:     d,
	}

	note, err := server.store.CreateKidNote(ctx, arg)

	if err != nil {
		errorsmanager.BuildError(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, note)
}

type updateKidNoteRequest struct {
	Note     string        `json:"note" binding:"min=0,max=250"`
	Presence []db.Presence `json:"presence" binding:"required,presence"`
	HasMeal  *bool         `json:"has_meal" binding:"required"`
}
type updateKidNoteParam struct {
	NoteID int64 `uri:"note_id" binding:"required"`
}

func (server *Server) updateKidNoteRoute(ctx *gin.Context) {
	var req updateKidNoteRequest
	var param updateKidNoteParam
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorsmanager.MapValidationErrors(ctx, err)
		return
	}
	if err := ctx.ShouldBindUri(&param); err != nil {
		errorsmanager.MapValidationErrors(ctx, err)
		return
	}
	kn, e := server.store.GetKidNote(ctx, param.NoteID)
	if e != nil {
		errorsmanager.BuildError(e, ctx)
		return
	}
	// User can only update note text
	arg := db.UpdateKidNoteParams{
		ID:       kn.ID,
		Note:     req.Note,
		Presence: req.Presence,
		HasMeal:  *req.HasMeal,
	}
	updatedNote, e := server.store.UpdateKidNote(ctx, arg)
	if e != nil {
		errorsmanager.BuildError(e, ctx)
		return
	}
	ctx.JSON(http.StatusOK, updatedNote)
}

type getKidNoteByPeriodRequest struct {
	Date1 string `json:"date1" binding:"required,date_format"`
	Date2 string `json:"date2" binding:"required,date_format"`
}

// Get note by period
func (server *Server) getKidNotesByPeriod(ctx *gin.Context) {
	var req getKidNoteByPeriodRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorsmanager.MapValidationErrors(ctx, err)
		return
	}
	dt1, _ := util.ConvertDate(req.Date1)
	dt2, _ := util.ConvertDate(req.Date2)
	arg := db.GetKidNotesByPeriodParams{
		Date:   dt1,
		Date_2: dt2,
	}
	tns, e := server.store.GetKidNotesByPeriod(ctx, arg)
	if e != nil {
		errorsmanager.BuildError(e, ctx)
		return
	}
	ctx.JSON(http.StatusOK, tns)
}

func (server *Server) getAllKidNotesRoute(ctx *gin.Context) {
	kns, e := server.store.GetAllKidNotes(ctx)
	if e != nil {
		errorsmanager.BuildError(e, ctx)
		return
	}
	ctx.JSON(http.StatusOK, kns)
}
