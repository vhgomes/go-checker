package server

import (
	"github.com/gin-gonic/gin"
	handlers2 "go-checker/internal/handlers"
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

	router.POST("/sites", siteHandler.CreateSite)
	router.GET("/sites", siteHandler.GetSites)
	router.GET("/sites/status/:siteId", siteHandler.GetAllSiteStatusBySiteId)
	router.GET("/sites/status/:siteId/:firstDate/:secondDate", siteHandler.GetAllSiteStatusBySiteIdAndDate)

	router.POST("/users", userHandler.RegisterUser)
	router.POST("/users", userHandler.Login)

	return router
}
