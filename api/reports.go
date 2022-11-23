package api

import (
	utils "birdie/api/utils"
	db "birdie/db/sqlc"
	"birdie/errorsmanager"
	"birdie/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GetAllReportInfoRequest struct {
	Date1 string `json:"date1" binding:"required,date_format"`
	Date2 string `json:"date2" binding:"required,date_format"`
}

type GetAllReportInfoResponse struct {
	KidNotes     []db.GetKidNotesByPeriodRow
	TeacherNotes []db.GetTeacherNotesByPeriodRow
	CarnetInfo   utils.KidCarnetInfo
}

func (server *Server) getAllReportInfo(ctx *gin.Context) {
	var req GetAllReportInfoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorsmanager.MapValidationErrors(ctx, err)
		return
	}
	rkn, e := getReportKidNotes(*server, ctx, req.Date1, req.Date2)
	if e != nil {
		errorsmanager.BuildError(e, ctx)
		return
	}
	rtn, e := getReportTeacherNotes(*server, ctx, req.Date1, req.Date2)
	if e != nil {
		errorsmanager.BuildError(e, ctx)
		return
	}
	rc, e := getReportCarnets(*server, ctx)
	if e != nil {
		errorsmanager.BuildError(e, ctx)
		return
	}
	resp := GetAllReportInfoResponse{
		KidNotes:     rkn,
		TeacherNotes: rtn,
		CarnetInfo:   rc,
	}
	ctx.JSON(http.StatusOK, resp)
}

func getReportKidNotes(server Server, ctx *gin.Context, date1 string, date2 string) ([]db.GetKidNotesByPeriodRow, error) {
	dt1, _ := util.ConvertDate(date1)
	dt2, _ := util.ConvertDate(date2)
	arg := db.GetKidNotesByPeriodParams{
		Date:   dt1,
		Date_2: dt2,
	}
	tns, e := server.store.GetKidNotesByPeriod(ctx, arg)
	if e != nil {
		return nil, e
	}
	return tns, nil
}

func getReportTeacherNotes(server Server, ctx *gin.Context, date1 string, date2 string) ([]db.GetTeacherNotesByPeriodRow, error) {
	dt1, _ := util.ConvertDate(date1)
	dt2, _ := util.ConvertDate(date2)
	arg := db.GetTeacherNotesByPeriodParams{
		Date:   dt1,
		Date_2: dt2,
	}
	tns, e := server.store.GetTeacherNotesByPeriod(ctx, arg)
	if e != nil {
		return nil, e
	}
	return tns, nil
}

func getReportCarnets(server Server, ctx *gin.Context) (utils.KidCarnetInfo, error) {
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
