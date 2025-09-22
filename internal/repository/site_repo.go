package repository

import "gorm.io/gorm"

type Site struct {
	ID            uint   `gorm:"primaryKey"`
	URL           string `gorm:"unique"`
	Status        string
	CheckInterval int
}

type SiteRepo struct {
	DB *gorm.DB
}

func NewSiteRepo(db *gorm.DB) *SiteRepo {
	return &SiteRepo{DB: db}
}

func (r *SiteRepo) AddSite(url string, interval int) error {
	if interval <= 0 {
		interval = 30 // default de 30s
	}

	site := Site{
		URL:           url,
		Status:        "unknown",
		CheckInterval: interval,
	}

	return r.DB.Create(&site).Error
}

func (r *SiteRepo) GetSites() ([]Site, error) {
	var sites []Site
	err := r.DB.Find(&sites).Error
	return sites, err
}

func (r *SiteRepo) UpdateStatus(id uint, status string) error {
	return r.DB.Model(&Site{}).
		Where("id = ?", id).
		Update("status", status).Error
}
