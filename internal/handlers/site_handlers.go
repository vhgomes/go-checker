package handlers

import (
	"github.com/gin-gonic/gin"
	"go-checker/internal/repository"
	"net/http"
)

type SiteHandler struct {
	siteRepo       repository.SiteRepo
	siteStatusRepo repository.SiteStatusRepo
}

func NewSiteHandler(siteRepo *repository.SiteRepo, siteStatusRepo *repository.SiteStatusRepo) *SiteHandler {
	return &SiteHandler{siteRepo: *siteRepo, siteStatusRepo: *siteStatusRepo}
}

func (h *SiteHandler) CreateSite(c *gin.Context) {
	var body struct {
		URL string `json:"url"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": err.Error(),
		})
		return
	}

	if err := h.siteRepo.AddSite(body.URL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  http.StatusInternalServerError,
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
	})
}

func (h *SiteHandler) GetSites(c *gin.Context) {
	sites, err := h.siteRepo.GetSites()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  http.StatusInternalServerError,
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": sites,
	})
}
