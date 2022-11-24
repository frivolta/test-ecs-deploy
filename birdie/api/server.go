package api

import (
	db "birdie/db/sqlc"
	"birdie/middleware"
	"birdie/util"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config util.Config
	store  db.Store
	Router *gin.Engine
}

func NewServer(store db.Store, config util.Config) (*Server, error) {
	server := &Server{
		config: config,
		store:  store,
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("date_format", validDay)
		v.RegisterValidation("presence", validPresence)
	}

	server.setupRouter(config)
	return server, nil
}
func CORSConfig(config util.Config) cors.Config {
	corsConfig := cors.DefaultConfig()
	if config.GinMode == "release" {
		fmt.Println("Cors in release mode")
		corsConfig.AllowOrigins = []string{"https://app.apecalendar.link", "https://www.app.apecalendar.link", "https://arachne.apecalendar.link", "https://www.arachne.apecalendar.link", "https://www.arachne.apecalendar.link/", "http://arachne-env.eu-west-1.elasticbeanstalk.com", "http://arachne-env.eu-west-1.elasticbeanstalk.com/", "https://arachne-env.eu-west-1.elasticbeanstalk.com/"}
	} else {
		fmt.Println("Cors in dev mode")
		corsConfig.AllowOrigins = []string{"*"}
	}
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Access-Control-Allow-Headers", "Access-Control-Allow-Origin", "access-control-allow-origin, access-control-allow-headers", "Content-Type", "X-XSRF-TOKEN", "Accept", "Origin", "X-Requested-With", "Authorization")
	corsConfig.AddAllowMethods("GET", "POST", "PUT", "DELETE", "OPTIONS")
	return corsConfig
}
func (server *Server) setupRouter(config util.Config) {
	router := gin.Default()
	router.Use(cors.New(CORSConfig(config)))

	//Setting auth middleware
	router.Use(func(context *gin.Context) {
		context.Set("firebaseAuth", server.config.AuthClient)
	})

	// Health
	router.GET("/", server.publicRoute)
	router.GET("/heartbeat", server.publicRoute)
	router.GET("/private", middleware.AuthMiddleware, middleware.DBUserMiddleware(&server.store), server.privateRoute)
	router.GET("/scoped", middleware.AuthMiddleware, middleware.DBUserMiddleware(&server.store), middleware.RoleAuthorizationMiddleware(db.RoleTEACHER), server.privateRoute)

	// Teachers
	teachersRoutes := router.Group("/api/v1/teachers", middleware.AuthMiddleware, middleware.DBUserMiddleware(&server.store), middleware.RoleAuthorizationMiddleware(db.RoleTEACHER))
	teachersRoutes.POST("/", server.createTeacherRoute)
	teachersRoutes.GET("/", server.getAllTeachersRoute)
	teachersRoutes.GET("/:id", server.getTeacherRoute)

	// Teacher notes
	teacherNotesRoutes := router.Group("/api/v1/teacher_notes", middleware.AuthMiddleware, middleware.DBUserMiddleware(&server.store), middleware.RoleAuthorizationMiddleware(db.RoleTEACHER))
	teacherNotesRoutes.POST("/:teacher_id", server.createTeacherNoteRoute)
	teacherNotesRoutes.POST("/date", server.getTeacherNotesByDate)
	teacherNotesRoutes.POST("/period", server.getTeacherNotesByPeriod)
	teacherNotesRoutes.GET("/", server.getAllTeacherNotesRoute)
	teacherNotesRoutes.PUT("/:note_id", server.updateTeacherNote)

	// Kids
	kidRoutes := router.Group("/api/v1/kids", middleware.AuthMiddleware, middleware.DBUserMiddleware(&server.store), middleware.RoleAuthorizationMiddleware(db.RoleTEACHER))
	kidRoutes.POST("/", server.createKidRoute)
	kidRoutes.GET("/:id", server.getKidRoute)
	kidRoutes.GET("/", server.getAllKidsRoute)

	// Kid Notes
	kidNotesRoutes := router.Group("/api/v1/kid_notes", middleware.AuthMiddleware, middleware.DBUserMiddleware(&server.store), middleware.RoleAuthorizationMiddleware(db.RoleTEACHER))
	kidNotesRoutes.GET("/", server.getAllKidNotesRoute)
	kidNotesRoutes.POST("/:kid_id", server.createKidNoteRoute)
	kidNotesRoutes.PUT("/:note_id", server.updateKidNoteRoute)
	kidNotesRoutes.POST("/period", server.getKidNotesByPeriod)

	// Carnets
	carnetRoutes := router.Group("/api/v1/carnets", middleware.AuthMiddleware, middleware.DBUserMiddleware(&server.store), middleware.RoleAuthorizationMiddleware(db.RoleTEACHER))
	carnetRoutes.POST("/:kid_id", server.createCarnetRoute)
	carnetRoutes.GET("/", server.getAllCarnetsRoute)
	carnetRoutes.PUT("/:carnet_id", server.updateCarnetRoute)

	// Report
	reportRoutes := router.Group("/api/v1/reports", middleware.AuthMiddleware, middleware.DBUserMiddleware(&server.store), middleware.RoleAuthorizationMiddleware(db.RoleADMIN))
	reportRoutes.POST("/monthly-report", server.getAllReportInfo)

	server.Router = router

}

func (server *Server) Start(address string) error {
	return server.Router.Run(address)
}
