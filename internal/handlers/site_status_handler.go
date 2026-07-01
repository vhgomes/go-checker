package handlers

import (
	"go-checker/internal/middlewares"
	"go-checker/internal/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type SiteStatusHandler struct {
	Repo *repository.SiteStatusRepo
}

func NewSiteStatusHandler(repo *repository.SiteStatusRepo) *SiteStatusHandler {
	return &SiteStatusHandler{Repo: repo}
}

func (h *SiteStatusHandler) GetAllSiteStatusBySiteIdPaginated(c *gin.Context) {
	userID, ok := middlewares.GetUserID(c)
	if !ok {
		return
	}

	siteId, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	result, err := h.Repo.GetAllSiteStatusBySiteIdPaginated(userID, uint(siteId), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *SiteStatusHandler) GetAllSiteStatusByDatePaginated(c *gin.Context) {
	userID, ok := middlewares.GetUserID(c)
	if !ok {
		return
	}

	siteId, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	from := c.Query("from")
	to := c.Query("to")

	firstDate, err := time.Parse("2006-01-02", from)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from date"})
		return
	}
	secondDate, err := time.Parse("2006-01-02", to)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to date"})
		return
	}

	result, err := h.Repo.GetAllSiteStatusBySiteIdAndDatePaginated(userID, uint(siteId), firstDate, secondDate, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *SiteStatusHandler) GetAllSiteStatusByStatusPaginated(c *gin.Context) {
	userID, ok := middlewares.GetUserID(c)
	if !ok {
		return
	}

	siteId, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	result, err := h.Repo.GetAllSiteStatusByStatusPaginated(userID, uint(siteId), status, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *SiteStatusHandler) GetLastSiteStatus(c *gin.Context) {
	userID, ok := middlewares.GetUserID(c)
	if !ok {
		return
	}

	siteId, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	result, err := h.Repo.GetLastSiteStatus(userID, uint(siteId))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *SiteStatusHandler) GetFirstSiteStatus(c *gin.Context) {
	userID, ok := middlewares.GetUserID(c)
	if !ok {
		return
	}

	siteId, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	result, err := h.Repo.GetFirstSiteStatus(userID, uint(siteId))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
