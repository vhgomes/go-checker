package repository

import (
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

func (r *SiteRepo) DeleteSite(id uint) error {
	return r.DB.Delete(&Site{ID: id}).Error
}
func (r *SiteRepo) GetSiteById(id uint) (error, error) {
	return r.DB.First(&Site{ID: id}).Error, nil
}
func (r *SiteRepo) GetSitesByUserId(userId uint) ([]Site, error) {
	var sites []Site
	err := r.DB.Find(&sites).Where("user_id = ?", userId).Error
	return sites, err
}

// Funções que estão fora do escopo atual
func (r *SiteRepo) UpdateStatus(id uint, status string) error {
	return r.DB.Model(&Site{}).
		Where("id = ?", id).
		Update("status", status).Error
}
func (r *SiteRepo) GetSites() ([]Site, error) {
	var sites []Site
	err := r.DB.Find(&sites).Error
	return sites, err
}
