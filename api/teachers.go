package api

import (
	db "birdie/db/sqlc"
	"birdie/errorsmanager"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

type teacherResponse struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

// Helper, generate Teacher response
func newTeacherResponse(t *db.Teacher) teacherResponse {
	return teacherResponse{
		ID:      t.ID,
		Name:    t.Name,
		Surname: t.Surname,
	}
}

type createTeacherRequest struct {
	Name    string `json:"name" binding:"required,min=1,max=50,alpha"`
	Surname string `json:"surname" binding:"required,min=1,max=50,alpha"`
}

func (server *Server) createTeacherRoute(ctx *gin.Context) {
	var req createTeacherRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorsmanager.ErrorResponse(err))
		return
	}
	arg := db.CreateTeacherParams{
		Name:    req.Name,
		Surname: req.Surname,
	}

	t, e := server.store.CreateTeacher(ctx, arg)
	if e != nil {
		ctx.JSON(http.StatusInternalServerError, errorsmanager.ErrorResponse(e))
	}
	ctx.JSON(http.StatusOK, newTeacherResponse(&t))
}

// Get all teachers ordered by name
func (server *Server) getAllTeachersRoute(ctx *gin.Context) {
	ts, e := server.store.GetAllTeachers(ctx)
	if e != nil {
		ctx.JSON(http.StatusInternalServerError, errorsmanager.ErrorResponse(e))
		return
	}
	ctx.JSON(http.StatusOK, ts)
}

type getTeacherRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// Get a single teacher
func (server *Server) getTeacherRoute(ctx *gin.Context) {
	var req getTeacherRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorsmanager.ErrorResponse(err))
		return
	}
	t, e := server.store.GetTeacher(ctx, req.ID)
	if e != nil {
		if e == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorsmanager.ErrorResponse(e))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorsmanager.ErrorResponse(e))
		return
	}

	ctx.JSON(http.StatusOK, newTeacherResponse(&t))
}
