package handlers

import (
	"context"
	"go-checker/internal/middlewares"
	"go-checker/internal/monitor"
	"go-checker/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SiteHandler struct {
	siteRepo       repository.SiteRepo
	siteStatusRepo repository.SiteStatusRepo
	monitorManager *monitor.MonitorManager
	ctx            context.Context
}

func NewSiteHandler(siteRepo *repository.SiteRepo, siteStatusRepo *repository.SiteStatusRepo, monitorManager *monitor.MonitorManager, ctx context.Context) *SiteHandler {
	return &SiteHandler{siteRepo: *siteRepo, siteStatusRepo: *siteStatusRepo, monitorManager: monitorManager, ctx: ctx}
}

func (h *SiteHandler) CreateSite(c *gin.Context) {
	userID, ok := middlewares.GetUserID(c)
	if !ok {
		return
	}

	var body struct {
		URL           string `json:"url"`
		CheckInterval int    `json:"check_interval"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "corpo da requisição invalido"})
		return
	}

	newSite, err := h.siteRepo.AddSite(body.URL, body.CheckInterval, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao criar o site"})
		return
	}

	h.monitorManager.Register(*newSite)

	c.JSON(http.StatusOK, gin.H{"message": "site criado com sucesso"})
}

func (h *SiteHandler) DeleteSite(c *gin.Context) {
	idParam := c.Param("id")
	siteID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id do site invalido"})
		return
	}

	userID, ok := middlewares.GetUserID(c)
	if !ok {
		return
	}

	if err := h.siteRepo.DeleteSite(uint(siteID), userID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "site deletado com sucesso"})
}

func (h *SiteHandler) UpdateSite(c *gin.Context) {
	idParam := c.Param("id")
	siteID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id do site invalido"})
		return
	}

	var body struct {
		URL           string `json:"url"`
		CheckInterval int    `json:"check_interval"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "corpo da requisição invalido"})
		return
	}

	userID, ok := middlewares.GetUserID(c)
	if !ok {
		return
	}

	if err := h.siteRepo.UpdateSite(uint(siteID), userID, body.URL, body.CheckInterval); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "site atualizado com sucesso"})
}

func (h *SiteHandler) GetSiteById(c *gin.Context) {
	idParam := c.Param("id")
	siteID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id do site invalido"})
		return
	}

	userID, ok := middlewares.GetUserID(c)
	if !ok {
		return
	}

	site, err := h.siteRepo.GetSiteById(uint(siteID), userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": site})
}

func (h *SiteHandler) GetAllSitesByUser(c *gin.Context) {
	userID, ok := middlewares.GetUserID(c)
	if !ok {
		return
	}

	sites, err := h.siteRepo.GetSitesByUserId(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha em pegar todos os sites"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": sites})
}

func (h *SiteHandler) GetDashboardByUser(c *gin.Context) {
	userID, ok := middlewares.GetUserID(c)
	if !ok {
		return
	}

	dashboard, err := h.siteRepo.GetAllSiteInfoByUserId(h.ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao juntar o dashboard"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": dashboard})
}
