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

	monitor.StartMonitoring(siteRepo)

	router.POST("/sites", siteHandler.CreateSite)
	router.GET("/sites", siteHandler.GetSites)

	return router
}
