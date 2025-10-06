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

	// repos
	siteRepo := repository2.NewSiteRepo(db)
	siteStatusRepo := repository2.NewSiteStatusRepo(db)
	userRepo := repository2.NewUserRepo(db)

	// handlers
	siteHandler := handlers2.NewSiteHandler(siteRepo, siteStatusRepo, ctx)
	userHandler := handlers2.NewUserHandler(userRepo)
	siteStatusHandler := handlers2.NewSiteStatusHandler(siteStatusRepo)

	// cronjobs
	monitor.StartMonitoring(ctx, siteRepo, siteStatusRepo)

	cronManager := cronjobs.NewJobManager(ctx)
	dashboardCronJob := cronjobs.NewDashboardCronJob(siteRepo, userRepo, redis, "*/30 * * * * *") // a cada 30s
	cronManager.RegisterJob(dashboardCronJob)
	cronManager.StartScheduler()

	auth := router.Group("/")
	auth.Use(middlewares.MiddlewareJWT())
	{
		// rotas basicas dos sites
		auth.POST("/sites", siteHandler.CreateSite)
		auth.GET("/sites", siteHandler.GetAllSitesByUser)
		auth.GET("/sites/:id", siteHandler.GetSiteById)
		auth.PUT("/sites/:id", siteHandler.UpdateSite)
		auth.DELETE("/sites/:id", siteHandler.DeleteSite)

		// rotas filtros
		status := auth.Group("/sites/:id/status")
		{
			status.GET("", siteStatusHandler.GetAllSiteStatusBySiteIdPaginated)
			status.GET("/date", siteStatusHandler.GetAllSiteStatusByDatePaginated)
			status.GET("/filter", siteStatusHandler.GetAllSiteStatusByStatusPaginated)
			status.GET("/last", siteStatusHandler.GetLastSiteStatus)
			status.GET("/first", siteStatusHandler.GetFirstSiteStatus)
		}
	}

	//rotas publicas
	router.POST("/register", userHandler.RegisterUser)
	router.POST("/login", userHandler.Login)

	return router
}
