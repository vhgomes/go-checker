package server

import (
	"github.com/gin-gonic/gin"
	"go-checker/internal/monitor"
	repository2 "go-checker/internal/repository"
	"gorm.io/gorm"
	"net/http"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	repository := repository2.NewSiteRepo(db)

	monitor.StartMonitoring(repository)

	router.POST("/sites", func(c *gin.Context) {
		var body struct {
			URL string `json:"url"`
		}

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":  http.StatusBadRequest,
				"error": err.Error(),
			})
		}

		if err := repository.AddSite(body.URL); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":  http.StatusInternalServerError,
				"error": err.Error(),
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
		})
	})

	router.GET("/sites", func(c *gin.Context) {
		sites, _ := repository.GetSites()
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": sites,
		})
	})

	return router
}
