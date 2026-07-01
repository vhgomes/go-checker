package server

import (
	"context"
	"go-checker/internal/cronjobs"
	handlers2 "go-checker/internal/handlers"
	"go-checker/internal/middlewares"
	"go-checker/internal/monitor"
	repository2 "go-checker/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupRouter(ctx context.Context, db *gorm.DB, redis *redis.Client) *gin.Engine {
	router := gin.Default()

	// repos
	siteRepo := repository2.NewSiteRepo(db)
	siteStatusRepo := repository2.NewSiteStatusRepo(db)
	userRepo := repository2.NewUserRepo(db)

	// monitor manager: inicia sites existentes e novos são colocados
	monitorManager := monitor.NewMonitorManager(ctx, siteRepo, siteStatusRepo)
	monitorManager.Start()

	// handlers
	siteHandler := handlers2.NewSiteHandler(siteRepo, siteStatusRepo, monitorManager)
	userHandler := handlers2.NewUserHandler(userRepo)
	siteStatusHandler := handlers2.NewSiteStatusHandler(siteStatusRepo)

	// cronjobs
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

		auth.GET("/dashboard", siteHandler.GetDashboardByUser)

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
	public := router.Group("/")
	public.Use(middlewares.AuthRateLimiter())
	{
		public.POST("/register", userHandler.RegisterUser)
		public.POST("/login", userHandler.Login)
	}

	return router
}
