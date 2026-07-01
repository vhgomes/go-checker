package monitor

import (
	"context"
	"go-checker/internal/repository"

	"go.uber.org/zap"
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
		zap.L().Error("Erro ao buscar sites para monitoramento", zap.Error(err))
		return
	}
	for _, site := range sites {
		go monitorSite(m.ctx, m.repo, m.statusRepo, site)
	}

	go m.listenForNewSites()
}

func (m *MonitorManager) Register(site repository.Site) {
	select {
	case m.newSites <- site:
		zap.L().Info("Site registrado para monitoramento dinâmico", zap.String("site_url", site.URL), zap.Uint("site_id", site.ID))
	default:
		zap.L().Warn("Canal de novos sites cheio, site não foi registrado", zap.String("site_url", site.URL), zap.Uint("site_id", site.ID))
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
