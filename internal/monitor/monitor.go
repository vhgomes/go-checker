package monitor

import (
	"go-checker/internal/repository"
	"log"
	"net/http"
	"time"
)

func StartMonitoring(repo *repository.SiteRepo, statusRepo *repository.SiteStatusRepo) {
	sites, _ := repo.GetSites()
	for _, site := range sites {
		go monitorSite(repo, statusRepo, site)
	}
}

func monitorSite(repo *repository.SiteRepo, statusRepo *repository.SiteStatusRepo, site repository.Site) {
	interval := 30 // default
	if site.CheckInterval > 0 {
		interval = site.CheckInterval
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		checkSite(repo, site, statusRepo)
	}
}

func checkSite(repo *repository.SiteRepo, site repository.Site, statusRepo *repository.SiteStatusRepo) {
	start := time.Now()
	resp, err := http.Get(site.URL)
	responseTime := time.Since(start).Seconds()

	status := "online"
	statusCode := 0
	if err != nil || resp.StatusCode >= 400 {
		status = "offline"
		if resp != nil {
			statusCode = resp.StatusCode
		}
	} else {
		statusCode = resp.StatusCode
	}

	if err := repo.UpdateStatus(site.ID, status); err != nil {
		log.Printf("Erro ao atualizar status do site %d: %v", site.ID, err)
	}

	if err := statusRepo.Insert(
		site.ID,
		status,
		statusCode,
		responseTime,
		time.Now(),
	); err != nil {
		log.Printf("Erro ao inserir histórico do site %d: %v", site.ID, err)
	}

	log.Printf("Site %s %s (statusCode=%d, responseTime=%.3fs)\n", site.URL, status, statusCode, responseTime)
}
