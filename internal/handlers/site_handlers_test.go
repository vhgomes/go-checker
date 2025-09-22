package handlers_test

import (
	"bytes"
	"encoding/json"
	"go-checker/internal/handlers"
	"go-checker/internal/repository"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*repository.SiteRepo, *repository.SiteStatusRepo) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal("falha ao criar DB:", err)
	}

	db.AutoMigrate(&repository.Site{}, &repository.SiteStatusHistory{})

	return repository.NewSiteRepo(db), repository.NewSiteStatusRepo(db)
}

func performRequest(handler gin.HandlerFunc, method, path string, body interface{}) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	var reqBody *bytes.Buffer
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBody)
	} else {
		reqBody = &bytes.Buffer{}
	}

	c.Request, _ = http.NewRequest(method, path, reqBody)
	c.Request.Header.Set("Content-Type", "application/json")
	handler(c)

	return w
}

func TestCreateSite(t *testing.T) {
	siteRepo, statusRepo := setupTestDB(t)
	h := handlers.NewSiteHandler(siteRepo, statusRepo)

	body := map[string]interface{}{
		"url":            "https://example.com",
		"check_interval": 60,
	}

	w := performRequest(h.CreateSite, "POST", "/sites", body)

	if w.Code != http.StatusOK {
		t.Fatalf("esperava status 200, recebeu %d", w.Code)
	}

	sites, _ := siteRepo.GetSites()
	if len(sites) != 1 {
		t.Fatal("site não foi criado no DB")
	}
	if sites[0].CheckInterval != 60 {
		t.Fatal("check_interval não foi salvo corretamente")
	}
}

func TestGetSites(t *testing.T) {
	siteRepo, statusRepo := setupTestDB(t)
	h := handlers.NewSiteHandler(siteRepo, statusRepo)

	_ = siteRepo.AddSite("https://example.com", 45)

	w := performRequest(h.GetSites, "GET", "/sites", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("esperava status 200, recebeu %d", w.Code)
	}
}

func TestInsertSiteStatus(t *testing.T) {
	siteRepo, statusRepo := setupTestDB(t)
	h := handlers.NewSiteHandler(siteRepo, statusRepo)

	siteRepo.AddSite("https://example.com", 30)

	body := map[string]interface{}{
		"site_id":       1,
		"status":        "online",
		"status_code":   200,
		"response_time": 0.123,
		"checked_at":    time.Now(),
	}

	w := performRequest(h.InsertSiteStatus, "POST", "/site-status", body)
	if w.Code != http.StatusCreated {
		t.Fatalf("esperava status 201, recebeu %d", w.Code)
	}
}

func TestGetAllSiteStatusBySiteId(t *testing.T) {
	siteRepo, statusRepo := setupTestDB(t)
	h := handlers.NewSiteHandler(siteRepo, statusRepo)

	siteRepo.AddSite("https://example.com", 30)
	statusRepo.Insert(1, "online", 200, 0.12, time.Now())

	w := performRequest(h.GetAllSiteStatusBySiteId, "GET", "/site-status/1", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("esperava status 200, recebeu %d", w.Code)
	}
}

func TestGetAllSiteStatusBySiteIdAndDate(t *testing.T) {
	siteRepo, statusRepo := setupTestDB(t)
	h := handlers.NewSiteHandler(siteRepo, statusRepo)

	siteRepo.AddSite("https://example.com", 30)
	now := time.Now()
	statusRepo.Insert(1, "online", 200, 0.12, now)

	// Datas no formato MM-DD-YYYY usado pelo handler
	path := "/site-status/1/" + now.Format("01-02-2006") + "/" + now.Format("01-02-2006")
	w := performRequest(h.GetAllSiteStatusBySiteIdAndDate, "GET", path, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("esperava status 200, recebeu %d", w.Code)
	}
}
