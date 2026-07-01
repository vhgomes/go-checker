package handlers

import (
	"fmt"
	"go-checker/internal/middlewares"
	"go-checker/internal/monitor"
	"go-checker/internal/repository"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type SiteHandler struct {
	siteRepo       repository.SiteRepo
	siteStatusRepo repository.SiteStatusRepo
	monitorManager *monitor.MonitorManager
}

func NewSiteHandler(siteRepo *repository.SiteRepo, siteStatusRepo *repository.SiteStatusRepo, monitorManager *monitor.MonitorManager) *SiteHandler {
	return &SiteHandler{siteRepo: *siteRepo, siteStatusRepo: *siteStatusRepo, monitorManager: monitorManager}
}

type siteRequest struct {
	URL           string `json:"url"            validate:"required,url"`
	CheckInterval int    `json:"check_interval" validate:"required,min=10"`
}

func validateSiteRequest(body siteRequest) error {
	if err := validate.Struct(body); err != nil {
		var ve validator.ValidationErrors
		if ok := (err.(validator.ValidationErrors)); ok != nil {
			ve = ok
		}
		if len(ve) > 0 {
			fe := ve[0]
			switch fe.Field() {
			case "URL":
				return fmt.Errorf("url inválida: deve ser uma URL completa")
			case "CheckInterval":
				return fmt.Errorf("check_interval deve ser no mínimo %s segundos", fe.Param())
			}
		}
		return err
	}
	// extra: garante que o scheme é http ou https
	u, err := url.ParseRequestURI(body.URL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		return fmt.Errorf("url inválida: use http:// ou https://")
	}
	return nil
}

func (h *SiteHandler) CreateSite(c *gin.Context) {
	ctx := c.Request.Context()
	userID, ok := middlewares.GetUserID(c)
	if !ok {
		return
	}

	var body siteRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "corpo da requisição invalido"})
		return
	}

	if err := validateSiteRequest(body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	newSite, err := h.siteRepo.AddSite(ctx, body.URL, body.CheckInterval, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao criar o site"})
		return
	}

	h.monitorManager.Register(*newSite)

	c.JSON(http.StatusOK, gin.H{"message": "site criado com sucesso"})
}

func (h *SiteHandler) DeleteSite(c *gin.Context) {
	ctx := c.Request.Context()
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

	if err := h.siteRepo.DeleteSite(ctx, uint(siteID), userID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "site deletado com sucesso"})
}

func (h *SiteHandler) UpdateSite(c *gin.Context) {
	ctx := c.Request.Context()
	idParam := c.Param("id")
	siteID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id do site invalido"})
		return
	}

	var body siteRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "corpo da requisição invalido"})
		return
	}

	if err := validateSiteRequest(body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	userID, ok := middlewares.GetUserID(c)
	if !ok {
		return
	}

	if err := h.siteRepo.UpdateSite(ctx, uint(siteID), userID, body.URL, body.CheckInterval); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "site atualizado com sucesso"})
}

func (h *SiteHandler) GetSiteById(c *gin.Context) {
	ctx := c.Request.Context()
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

	site, err := h.siteRepo.GetSiteById(ctx, uint(siteID), userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": site})
}

func (h *SiteHandler) GetAllSitesByUser(c *gin.Context) {
	ctx := c.Request.Context()
	userID, ok := middlewares.GetUserID(c)
	if !ok {
		return
	}

	sites, err := h.siteRepo.GetSitesByUserId(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha em pegar todos os sites"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": sites})
}

func (h *SiteHandler) GetDashboardByUser(c *gin.Context) {
	ctx := c.Request.Context()
	userID, ok := middlewares.GetUserID(c)
	if !ok {
		return
	}

	dashboard, err := h.siteRepo.GetAllSiteInfoByUserId(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao juntar o dashboard"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": dashboard})
}
