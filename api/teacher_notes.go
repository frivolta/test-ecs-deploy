package api

import (
	db "birdie/db/sqlc"
	"birdie/errorsmanager"
	"birdie/util"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

type createTeacherNoteRequest struct {
	Note string `json:"note" binding:"required,min=1,max=250"`
	Date string `json:"date" binding:"required,date_format"`
}
type createTeacherNoteParams struct {
	TeacherID int64 `uri:"teacher_id" binding:"required,min=1"`
}

func (server *Server) createTeacherNoteRoute(ctx *gin.Context) {
	var par createTeacherNoteParams
	var req createTeacherNoteRequest
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

	// Check if teacher exists
	tc, err := server.store.GetTeacher(ctx, par.TeacherID)
	if err != nil {
		errorsmanager.BuildError(err, ctx)
		return
	}
	convertedDate, _ := util.ConvertDate(req.Date)

	// Create args
	arg := db.CreateTeacherNoteParams{
		Note: req.Note,
		TeacherID: sql.NullInt64{
			Int64: tc.ID,
			Valid: true,
		},
		Date: convertedDate,
	}
	note, err := server.store.CreateTeacherNote(ctx, arg)
	if err != nil {
		errorsmanager.BuildError(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, note)
}

type updateTeacherNoteRequest struct {
	Note string `json:"note" binding:"required,min=1,max=250"`
}
type updateTeacherNoteParam struct {
	NoteID int64 `uri:"note_id" binding:"required"`
}

func (server *Server) updateTeacherNote(ctx *gin.Context) {
	var req updateTeacherNoteRequest
	var param updateTeacherNoteParam
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorsmanager.MapValidationErrors(ctx, err)
		return
	}
	if err := ctx.ShouldBindUri(&param); err != nil {
		errorsmanager.MapValidationErrors(ctx, err)
		return
	}
	tn, e := server.store.GetTeacherNote(ctx, param.NoteID)
	if e != nil {
		errorsmanager.BuildError(e, ctx)
		return
	}
	// User can only update note text
	arg := db.UpdateNoteParams{
		ID:   tn.ID,
		Note: req.Note,
		Date: tn.Date,
	}
	updatedNote, e := server.store.UpdateNote(ctx, arg)
	if e != nil {
		errorsmanager.BuildError(e, ctx)
		return
	}
	ctx.JSON(http.StatusOK, updatedNote)
}

type getTeacherNoteByDateRequest struct {
	Date string `json:"date" binding:"required,date_format"`
}

// Get note by date
func (server *Server) getTeacherNotesByDate(ctx *gin.Context) {
	var req getTeacherNoteByDateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorsmanager.MapValidationErrors(ctx, err)
		return
	}
	dt, _ := util.ConvertDate(req.Date)
	tn, e := server.store.GetTeacherNotesByDate(ctx, dt)
	if e != nil {
		errorsmanager.BuildError(e, ctx)
		return
	}
	ctx.JSON(http.StatusOK, tn)
}

type getTeacherNoteByPeriodRequest struct {
	Date1 string `json:"date1" binding:"required,date_format"`
	Date2 string `json:"date2" binding:"required,date_format"`
}

// Get note by period
func (server *Server) getTeacherNotesByPeriod(ctx *gin.Context) {
	var req getTeacherNoteByPeriodRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorsmanager.MapValidationErrors(ctx, err)
		return
	}
	dt1, _ := util.ConvertDate(req.Date1)
	dt2, _ := util.ConvertDate(req.Date2)
	arg := db.GetTeacherNotesByPeriodParams{
		Date:   dt1,
		Date_2: dt2,
	}
	tns, e := server.store.GetTeacherNotesByPeriod(ctx, arg)
	if e != nil {
		errorsmanager.BuildError(e, ctx)
		return
	}
	ctx.JSON(http.StatusOK, tns)
}

// Get all notes
func (server *Server) getAllTeacherNotesRoute(ctx *gin.Context) {
	tns, e := server.store.GetAllTeacherNotes(ctx)
	if e != nil {
		errorsmanager.BuildError(e, ctx)
		return
	}
	ctx.JSON(http.StatusOK, tns)
}
