package server

import (
	"github.com/gin-gonic/gin"
	handlers2 "go-checker/internal/handlers"
	"go-checker/internal/middlewares"
	"go-checker/internal/monitor"
	repository2 "go-checker/internal/repository"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	siteRepo := repository2.NewSiteRepo(db)
	siteStatusRepo := repository2.NewSiteStatusRepo(db)
	siteHandler := handlers2.NewSiteHandler(siteRepo, siteStatusRepo)
	userRepo := repository2.NewUserRepo(db)
	userHandler := handlers2.NewUserHandler(userRepo)

	monitor.StartMonitoring(siteRepo, siteStatusRepo)

	authProtected := router.Group("/")
	authProtected.Use(middlewares.MiddlewareJWT())
	{
		authProtected.POST("/sites", siteHandler.CreateSite)
		authProtected.GET("/sites", siteHandler.GetSites)
		//authProtected.GET("/sites/status/:siteId", siteHandler.GetAllSiteStatusBySiteId)
		//authProtected.GET("/sites/status/:siteId/:firstDate/:secondDate", siteHandler.GetAllSiteStatusBySiteIdAndDate)
	}

	router.POST("/register", userHandler.RegisterUser)
	router.POST("/login", userHandler.Login)

	return router
}
