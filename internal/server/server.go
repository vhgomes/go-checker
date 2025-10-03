package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go-checker/internal/cronjobs"
	handlers2 "go-checker/internal/handlers"
	"go-checker/internal/middlewares"
	"go-checker/internal/monitor"
	repository2 "go-checker/internal/repository"
	"gorm.io/gorm"
)

func SetupRouter(ctx context.Context, db *gorm.DB, redis *redis.Client) *gin.Engine {
	router := gin.Default()

	siteRepo := repository2.NewSiteRepo(db)
	siteStatusRepo := repository2.NewSiteStatusRepo(db)
	userRepo := repository2.NewUserRepo(db)

	siteHandler := handlers2.NewSiteHandler(siteRepo, siteStatusRepo)
	userHandler := handlers2.NewUserHandler(userRepo)

	monitor.StartMonitoring(ctx, siteRepo, siteStatusRepo)

	cronManager := cronjobs.NewJobManager(ctx)
	dashboardCronJob := cronjobs.NewDashboardCronJob(siteRepo, userRepo, redis, "*/30 * * * * *") // a cada 30s
	cronManager.RegisterJob(dashboardCronJob)
	cronManager.StartScheduler()

	auth := router.Group("/")
	auth.Use(middlewares.MiddlewareJWT())
	{
		auth.POST("/sites", siteHandler.CreateSite)
		auth.GET("/sites", siteHandler.GetAllSitesByUser)
		auth.GET("/sites/:id", siteHandler.GetSiteById)
		auth.PUT("/sites/:id", siteHandler.UpdateSite)
		auth.DELETE("/sites/:id", siteHandler.DeleteSite)
	}

	router.POST("/register", userHandler.RegisterUser)
	router.POST("/login", userHandler.Login)

	return router
}
