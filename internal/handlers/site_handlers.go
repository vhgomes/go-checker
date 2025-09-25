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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid site id"})
		return
	}

	userAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}
	userID := uint(userAny.(float64))

	site, err := h.siteRepo.GetSiteById(uint(siteID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "site not found"})
		return
	}

	if site.UserId != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you are not the owner of this site"})
		return
	}

	if err := h.siteRepo.DeleteSite(uint(siteID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete site"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "site deleted successfully"})
}

func (h *SiteHandler) UpdateSite(c *gin.Context) {
	idParam := c.Param("id")
	siteID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid site id"})
		return
	}

	var body struct {
		URL           string `json:"url"`
		CheckInterval int    `json:"check_interval"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	userAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}
	userID := uint(userAny.(float64))

	site, err := h.siteRepo.GetSiteById(uint(siteID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "site not found"})
		return
	}

	if site.UserId != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you are not the owner of this site"})
		return
	}

	if err := h.siteRepo.UpdateSite(uint(siteID), body.URL, body.CheckInterval); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update site"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "site updated successfully"})
}

func (h *SiteHandler) GetSiteById(c *gin.Context) {
	idParam := c.Param("id")
	siteID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid site id"})
		return
	}

	userAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}
	userID := uint(userAny.(float64))

	site, err := h.siteRepo.GetSiteById(uint(siteID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "site not found"})
		return
	}

	if site.UserId != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you are not the owner of this site"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": site})
}

func (h *SiteHandler) GetAllSitesByUser(c *gin.Context) {
	sites, err := h.siteRepo.GetSites()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  http.StatusInternalServerError,
			"error": "failed to get all sites",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": sites,
	})
}

// Funções fora do escopo atual de desenvolvimento
func (h *SiteHandler) GetAllSiteStatusByUser(c *gin.Context) {
	userAny, exists := c.Get("user_id")
	pageParam := c.DefaultQuery("page", "1")
	pageSizeParam := c.DefaultQuery("page_size", "10")

	page, _ := strconv.Atoi(pageParam)
	pageSize, _ := strconv.Atoi(pageSizeParam)

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": "bad request, failed to get user id",
		})
		return
	}

	userFloat, ok := userAny.(float64)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": "bad request, failed to parse user id",
		})
		return
	}

	userId := uint(userFloat)

	sites, err := h.siteRepo.GetSitesByUserId(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  http.StatusInternalServerError,
			"error": "failed to get sites by user",
		})
		return
	}

	if len(sites) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code":  http.StatusNotFound,
			"error": "no sites found for this user",
		})
		return
	}

	allStatus := make(map[uint][]repository.SiteStatusHistory)
	for _, site := range sites {
		status, err := h.siteStatusRepo.GetAllSiteStatusBySiteIdPaginated(site.ID, page, pageSize)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":  http.StatusInternalServerError,
				"error": "error getting site status",
			})
			return
		}
		allStatus[site.ID] = status
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": allStatus,
	})
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
