package handlers

import (
	"github.com/gin-gonic/gin"
	"go-checker/internal/repository"
	"net/http"
	"strconv"
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
		URL           string `json:"url"`
		CheckInterval int    `json:"check_interval"`
	}

	userAny, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "não autorizado"})
		return
	}

	userID := uint(userAny.(float64))

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "corpo da requisição invalido"})
		return
	}

	if err := h.siteRepo.AddSite(body.URL, body.CheckInterval, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao criar o site"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "site criado com sucesso"})
}

func (h *SiteHandler) DeleteSite(c *gin.Context) {
	idParam := c.Param("id")
	siteID, err := strconv.ParseUint(idParam, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id do site invalido"})
		return
	}

	userAny, _ := c.Get("user_id")
	userID := uint(userAny.(float64))

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

	userAny, _ := c.Get("user_id")
	userID := uint(userAny.(float64))

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

	userAny, _ := c.Get("user_id")
	userID := uint(userAny.(float64))

	site, err := h.siteRepo.GetSiteById(uint(siteID), userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": site})
}

func (h *SiteHandler) GetAllSitesByUser(c *gin.Context) {
	userAny, _ := c.Get("user_id")
	userID := uint(userAny.(float64))

	sites, err := h.siteRepo.GetSitesByUserId(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha em pegar todos os sites"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": sites})
}
