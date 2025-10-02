package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
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

type SiteRepo struct {
	DB *gorm.DB
}

func NewSiteRepo(db *gorm.DB) *SiteRepo {
	return &SiteRepo{DB: db}
}

func (r *SiteRepo) AddSite(url string, interval int, user uint) error {
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

	return r.DB.Create(&site).Error
}

func (r *SiteRepo) UpdateSite(siteId uint, userId uint, url string, interval int) error {
	var site Site
	if err := r.DB.Where("id = ? AND user_id = ?", siteId, userId).First(&site).Error; err != nil {
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

func (r *SiteRepo) DeleteSite(siteId uint, userId uint) error {
	result := r.DB.Where("id = ? AND user_id = ?", siteId, userId).Delete(&Site{})

	if result.RowsAffected == 0 {
		return errors.New("site não encontrado ou site não pertence ao user")
	}

	return result.Error
}

func (r *SiteRepo) GetSiteById(siteId uint, userId uint) (*Site, error) {
	var site Site

	if err := r.DB.Where("id = ? AND user_id = ?", siteId, userId).First(&site).Error; err != nil {
		return nil, err
	}

	return &site, nil
}

func (r *SiteRepo) GetSitesByUserId(userId uint) ([]Site, error) {
	var sites []Site
	err := r.DB.Find(&sites).Where("user_id = ?", userId).Error
	return sites, err
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
