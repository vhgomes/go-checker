package repository

import (
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

func (r *SiteRepo) GetSitesByUserId(userId uint) ([]Site, error) {
	var sites []Site
	err := r.DB.Find(&sites).Where("user_id = ?", userId).Error
	return sites, err
}

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
