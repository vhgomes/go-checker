package monitor

import (
	"context"
	"go-checker/internal/repository"
	"net/http"
	"time"

	"go.uber.org/zap"
)

const (
	maxRetries     = 3
	requestTimeout = 10 * time.Second
	baseBackoff    = 500 * time.Millisecond
)

var httpClient = &http.Client{
	Timeout: requestTimeout,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
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
			zap.L().Info("Monitoramento encerrado para site", zap.String("site_url", site.URL))
			return
		case <-ticker.C:
			checkSite(ctx, repo, statusRepo, site)
		}
	}
}

func checkSite(ctx context.Context, repo *repository.SiteRepo, statusRepo *repository.SiteStatusRepo, site repository.Site) {
	checkCtx, cancel := context.WithTimeout(ctx, requestTimeout*time.Duration(maxRetries)+baseBackoff*time.Duration(maxRetries))
	defer cancel()

	statusCode, responseTime := doRequestWithRetry(checkCtx, site.URL)

	status := "online"
	if statusCode == 0 || statusCode >= 400 {
		status = "offline"
	}

	if err := repo.UpdateStatus(checkCtx, site.ID, status); err != nil {
		zap.L().Error("Erro ao atualizar status do site", zap.Uint("site_id", site.ID), zap.Error(err))
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
		zap.L().Error("Erro ao inserir histórico do site", zap.Uint("site_id", site.ID), zap.Error(err))
		return
	}

	zap.L().Info("Site verificado",
		zap.String("site_url", site.URL),
		zap.String("status", status),
		zap.Int("status_code", statusCode),
		zap.Float64("response_time", responseTime),
	)
}

func doRequestWithRetry(ctx context.Context, url string) (int, float64) {
	var lastStatusCode int
	var lastResponseTime float64

	for attempt := range maxRetries {
		if attempt > 0 {
			backoff := baseBackoff * (1 << attempt) // 1s, 2s, 4s...
			select {
			case <-ctx.Done():
				return lastStatusCode, lastResponseTime
			case <-time.After(backoff):
			}
		}

		statusCode, responseTime, err := doRequest(ctx, url)
		lastStatusCode = statusCode
		lastResponseTime = responseTime

		if err != nil {
			zap.L().Warn("Tentativa falhou", zap.Int("attempt", attempt+1), zap.Int("max_retries", maxRetries), zap.String("url", url), zap.Error(err))
			continue
		}

		return statusCode, responseTime
	}

	return lastStatusCode, lastResponseTime
}

func doRequest(ctx context.Context, url string) (statusCode int, responseTime float64, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		return 0, 0, err
	}

	req.Header.Set("User-Agent", "go-checker/1.0")

	start := time.Now()
	resp, err := httpClient.Do(req)
	responseTime = time.Since(start).Seconds()

	if err != nil {
		req2, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		req2.Header.Set("User-Agent", "go-checker/1.0")

		start = time.Now()
		resp, err = httpClient.Do(req2)
		responseTime = time.Since(start).Seconds()
		if err != nil {
			return 0, responseTime, err
		}
	}

	defer resp.Body.Close()
	return resp.StatusCode, responseTime, nil
}
