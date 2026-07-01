package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Site struct {
	ID            uint `gorm:"primaryKey"`
	URL           string
	UserId        uint
	Status        string
	CheckInterval int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type InfoDashboardUser struct {
	TotalSites           int64
	NumberOfSitesOnline  int64
	NumberOfSitesOffline int64
	ResponseTimeAvg      float64
	LastEvents           []SiteStatusHistory
}

type SiteRepo struct {
	DB *gorm.DB
}

func NewSiteRepo(db *gorm.DB) *SiteRepo {
	return &SiteRepo{DB: db}
}

func (r *SiteRepo) AddSite(ctx context.Context, url string, interval int, user uint) (*Site, error) {
	if interval <= 0 {
		interval = 60
	}

	site := Site{
		URL:           url,
		Status:        "unknown",
		CheckInterval: interval,
		UserId:        user,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := r.DB.WithContext(ctx).Create(&site).Error; err != nil {
		return nil, err
	}

	return &site, nil
}

func (r *SiteRepo) UpdateSite(ctx context.Context, siteId uint, userId uint, url string, interval int) error {
	var site Site
	if err := r.DB.WithContext(ctx).Where("id = ? AND user_id = ?", siteId, userId).First(&site).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("site não encontrado ou site não pertence ao user")
		}
		return err
	}

	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if url != "" {
		updates["url"] = url
	}

	if interval > 0 {
		updates["check_interval"] = interval
	}

	return r.DB.Model(&site).Updates(updates).Error
}

func (r *SiteRepo) DeleteSite(ctx context.Context, siteId uint, userId uint) error {
	result := r.DB.WithContext(ctx).Where("id = ? AND user_id = ?", siteId, userId).Delete(&Site{})

	if result.RowsAffected == 0 {
		return errors.New("site não encontrado ou site não pertence ao user")
	}

	return result.Error
}

func (r *SiteRepo) GetSiteById(ctx context.Context, siteId uint, userId uint) (*Site, error) {
	var site Site

	if err := r.DB.WithContext(ctx).Where("id = ? AND user_id = ?", siteId, userId).First(&site).Error; err != nil {
		return nil, err
	}

	return &site, nil
}

func (r *SiteRepo) GetSitesByUserId(ctx context.Context, userId uint) ([]Site, error) {
	var sites []Site
	err := r.DB.WithContext(ctx).Where("user_id = ?", userId).Find(&sites).Error
	return sites, err
}

// Funções para filtros //

func (r *SiteRepo) GetSitesBySiteStatus(ctx context.Context, siteStatus string, userId uint) ([]Site, error) {
	var sites []Site
	err := r.DB.WithContext(ctx).Where("status = ? AND user_id = ?", siteStatus, userId).Find(&sites).Error
	return sites, err
}

func (r *SiteRepo) GetAllSiteInfoByUserId(ctx context.Context, userId uint) (InfoDashboardUser, error) {
	var numberOfOnlineSites int64
	var numberOfOfflineSites int64
	var totalSites int64
	var responseTimeAvg float64
	var lastEvents []SiteStatusHistory

	if err := r.DB.WithContext(ctx).
		Model(&Site{}).
		Where("user_id = ? AND status = ?", userId, "online").
		Count(&numberOfOnlineSites).Error; err != nil {
		return InfoDashboardUser{}, err
	}

	if err := r.DB.WithContext(ctx).
		Model(&Site{}).
		Where("user_id = ? AND status = ?", userId, "offline").
		Count(&numberOfOfflineSites).Error; err != nil {
		return InfoDashboardUser{}, err
	}

	if err := r.DB.WithContext(ctx).
		Model(&Site{}).
		Where("user_id = ?", userId).
		Count(&totalSites).Error; err != nil {
		return InfoDashboardUser{}, err
	}

	if err := r.DB.WithContext(ctx).
		Model(&SiteStatusHistory{}).
		Select("AVG(response_time)").
		Joins("INNER JOIN sites ON sites.id = site_status_histories.site_id").
		Where("sites.user_id = ?", userId).
		Scan(&responseTimeAvg).Error; err != nil {
		return InfoDashboardUser{}, err
	}

	// está ordenando por checked_at porém apos a modificação dos campos sera colocado created_at
	if err := r.DB.WithContext(ctx).
		Joins("INNER JOIN sites ON sites.id = site_status_histories.site_id").
		Where("sites.user_id = ?", userId).
		Order("site_status_histories.checked_at DESC").
		Limit(10).
		Find(&lastEvents).Error; err != nil {
		return InfoDashboardUser{}, err
	}

	response := InfoDashboardUser{
		TotalSites:           totalSites,
		NumberOfSitesOnline:  numberOfOnlineSites,
		NumberOfSitesOffline: numberOfOfflineSites,
		ResponseTimeAvg:      responseTimeAvg,
		LastEvents:           lastEvents,
	}

	return response, nil
}

// Funções de Monitoramento dos sites //
func (r *SiteRepo) UpdateStatus(ctx context.Context, id uint, status string) error {
	return r.DB.WithContext(ctx).
		Model(&Site{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *SiteRepo) GetAllSitesToMonitoring(ctx context.Context) ([]Site, error) {
	var sites []Site
	err := r.DB.WithContext(ctx).Find(&sites).Error
	return sites, err
}
