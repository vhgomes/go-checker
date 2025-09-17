package monitor

import (
	"go-checker/internal/repository"
	"log"
	"net/http"
	"time"
)

func StartMonitoring(repo *repository.SiteRepo) {
	ticker := time.NewTicker(30 * time.Second) // checa a cada 30s
	go func() {
		for {
			select {
			case <-ticker.C:
				sites, _ := repo.GetSites()
				for _, site := range sites {
					go checkSite(repo, site)
				}
			}
		}
	}()
}

func checkSite(repo *repository.SiteRepo, site repository.Site) {
	resp, err := http.Get(site.URL)
	if err != nil || resp.StatusCode >= 400 {
		log.Printf("Site %s OFFLINE\n", site.URL)
		err := repo.UpdateStatus(site.ID, "offline")
		if err != nil {
			log.Fatal("Erro ao dar update no site:", site.ID)
		}
		return
	}
	log.Printf("Site %s ONLINE\n", site.URL)
	err = repo.UpdateStatus(site.ID, "online")
	if err != nil {
		log.Fatal("Erro ao dar update no site:", site.ID)
	}
}
