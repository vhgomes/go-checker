package monitor

import (
	"context"
	"go-checker/internal/repository"
	"log"
	"math/rand"
	"time"
)

func StartMonitoring(ctx context.Context, repo *repository.SiteRepo, statusRepo *repository.SiteStatusRepo) {
	sites, err := repo.GetAllSitesToMonitoring(ctx)

	if err != nil {
		log.Println("❌ Erro ao buscar sites para monitoramento:", err)
		return
	}

	for _, site := range sites {
		go monitorSite(ctx, repo, statusRepo, site)
	}
}

func monitorSite(ctx context.Context, repo *repository.SiteRepo, statusRepo *repository.SiteStatusRepo, site repository.Site) {
	interval := 30
	if site.CheckInterval > 0 {
		interval = site.CheckInterval
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("🛑 Monitoramento encerrado para site %s\n", site.URL)
			return
		case <-ticker.C:
			checkSiteRandom(ctx, repo, site, statusRepo)
		}
	}
}

func checkSiteRandom(ctx context.Context, repo *repository.SiteRepo, site repository.Site, statusRepo *repository.SiteStatusRepo) {
	checkCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	responseTime := rand.Float64()*1.9 + 0.1

	statusCodes := []int{200, 200, 200, 404, 500} // mais chances de 200
	statusCode := statusCodes[rand.Intn(len(statusCodes))]

	status := "online"
	if statusCode >= 400 {
		status = "offline"
	}

	if err := repo.UpdateStatus(checkCtx, site.ID, status); err != nil {
		log.Printf("❌ Erro ao atualizar status do site %d: %v", site.ID, err)
		return
	}

	if err := statusRepo.Insert(
		checkCtx,
		site.ID,
		status,
		statusCode,
		responseTime,
		time.Now(),
	); err != nil {
		log.Printf("❌ Erro ao inserir histórico do site %d: %v", site.ID, err)
		return
	}

	log.Printf("✅ Site %s %s (statusCode=%d, responseTime=%.3fs)\n",
		site.URL, status, statusCode, responseTime)
}

//func checkSite(repo *repository.SiteRepo, site repository.Site, statusRepo *repository.SiteStatusRepo) {
//	start := time.Now()
//	resp, err := http.Get(site.URL)
//	responseTime := time.Since(start).Seconds()
//
//	status := "online"
//	statusCode := 0
//	if err != nil || resp.StatusCode >= 400 {
//		status = "offline"
//		if resp != nil {
//			statusCode = resp.StatusCode
//		}
//	} else {
//		statusCode = resp.StatusCode
//	}
//
//	if err := repo.UpdateStatus(site.ID, status); err != nil {
//		log.Printf("Erro ao atualizar status do site %d: %v", site.ID, err)
//	}
//
//	if err := statusRepo.Insert(
//		site.ID,
//		status,
//		statusCode,
//		responseTime,
//		time.Now(),
//	); err != nil {
//		log.Printf("Erro ao inserir histórico do site %d: %v", site.ID, err)
//	}
//
//	log.Printf("Site %s %s (statusCode=%d, responseTime=%.3fs)\n", site.URL, status, statusCode, responseTime)
//}
