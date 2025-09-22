package handlers

import (
	"github.com/gin-gonic/gin"
	"go-checker/internal/repository"
	"log"
	"net/http"
	"strconv"
	"time"
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
		CheckInterval int    `json:"check_interval"` // novo campo
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": "bad request, failed to bindjson",
		})
		return
	}

	// agora passa URL + intervalo
	if err := h.siteRepo.AddSite(body.URL, body.CheckInterval); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  http.StatusInternalServerError,
			"error": "error adding site",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "site created successfully",
	})
}

func (h *SiteHandler) GetSites(c *gin.Context) {
	sites, err := h.siteRepo.GetSites()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  http.StatusInternalServerError,
			"error": "failed to get all sites",
		})
		return
	}

	// já retorna o check_interval no JSON
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": sites,
	})
}

func (h *SiteHandler) InsertSiteStatus(c *gin.Context) {
	var body struct {
		SiteID       uint      `json:"site_id"`
		Status       string    `json:"status"`
		StatusCode   int       `json:"status_code"`
		ResponseTime float64   `json:"response_time"`
		CheckedAt    time.Time `json:"checked_at"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": "bad request, failed to bindjson",
		})
		return
	}

	if err := h.siteStatusRepo.Insert(
		body.SiteID,
		body.Status,
		body.StatusCode,
		body.ResponseTime,
		body.CheckedAt,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  http.StatusInternalServerError,
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code": http.StatusCreated,
		"data": "site status inserido com sucesso",
	})

	return
}

func (h *SiteHandler) GetAllSiteStatusBySiteId(c *gin.Context) {
	idParam := c.Param("siteId")

	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": "bad request, failed to get site id",
		})
		return
	}

	siteId, err := strconv.ParseUint(idParam, 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": "bad request, failed to parse site id",
		})
		return
	}

	status, err := h.siteStatusRepo.GetAllSiteStatusBySiteId(uint(siteId))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  http.StatusInternalServerError,
			"error": err.Error(),
		})
		return
	}

	if status == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":  http.StatusNotFound,
			"error": "no status registered on this site",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": status,
	})

	return
}

func (h *SiteHandler) GetAllSiteStatusBySiteIdAndDate(c *gin.Context) {
	idParam := c.Param("siteId")
	firstDateStr := c.Param("firstDate")
	secondDateStr := c.Param("secondDate")

	layout := "01-02-2006"

	firstDate, err := time.Parse(layout, firstDateStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": "bad request, failed to parse first date",
		})
		return
	}

	secondDate, err := time.Parse(layout, secondDateStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": "bad request, failed to parse second date",
		})
		return
	}

	if firstDate.After(secondDate) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": "bad request, first date is after second date",
		})
		return
	}

	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": "bad request, failed to get siteid",
		})
		return
	}

	siteId, err := strconv.ParseUint(idParam, 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": "bad request, failed to parse siteid",
		})
		return
	}

	status, err := h.siteStatusRepo.GetAllSiteStatusBySiteIdAndDate(uint(siteId), firstDate, secondDate)
	log.Println(err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  http.StatusInternalServerError,
			"error": "error getting site status",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": status,
	})

	return
}
