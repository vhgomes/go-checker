package monitor

import (
	"context"
	"go-checker/internal/repository"
	"log"
)

type MonitorManager struct {
	repo       *repository.SiteRepo
	statusRepo *repository.SiteStatusRepo
	newSites   chan repository.Site
	ctx        context.Context
}

func NewMonitorManager(ctx context.Context, repo *repository.SiteRepo, statusRepo *repository.SiteStatusRepo) *MonitorManager {
	return &MonitorManager{
		repo:       repo,
		statusRepo: statusRepo,
		newSites:   make(chan repository.Site, 10),
		ctx:        ctx,
	}
}

func (m *MonitorManager) Start() {
	sites, err := m.repo.GetAllSitesToMonitoring(m.ctx)
	if err != nil {
		log.Println("❌ Erro ao buscar sites para monitoramento:", err)
	} else {
		for _, site := range sites {
			go monitorSite(m.ctx, m.repo, m.statusRepo, site)
		}
	}

	go m.listenForNewSites()
}

func (m *MonitorManager) Register(site repository.Site) {
	select {
	case m.newSites <- site:
		log.Printf("📡 Site %s registrado para monitoramento dinâmico\n", site.URL)
	default:
		log.Printf("⚠️ Canal de novos sites cheio, site %s não foi registrado\n", site.URL)
	}
}

func (m *MonitorManager) listenForNewSites() {
	for {
		select {
		case <-m.ctx.Done():
			return
		case site := <-m.newSites:
			go monitorSite(m.ctx, m.repo, m.statusRepo, site)
		}
	}
}
