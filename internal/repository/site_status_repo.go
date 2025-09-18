package repository

import (
	"gorm.io/gorm"
	"time"
)

type SiteStatusHistory struct {
	ID           uint `gorm:"primaryKey"`
	SiteID       uint `gorm:"index"`
	Status       string
	StatusCode   *int
	ResponseTime float64
	CheckedAt    time.Time
}

type SiteStatusRepo struct {
	DB *gorm.DB
}

func NewSiteStatusRepo(db *gorm.DB) *SiteStatusRepo {
	return &SiteStatusRepo{
		DB: db,
	}
}

func (nsr *SiteStatusRepo) Insert(siteId uint, status string, statusCode int, responseTime float64, checked time.Time) error {
	return nsr.DB.Create(&SiteStatusHistory{
		SiteID:       siteId,
		Status:       status,
		StatusCode:   &statusCode,
		ResponseTime: responseTime,
		CheckedAt:    checked,
	}).Error
}

func (nsr *SiteStatusRepo) GetAllSiteStatusBySiteId(siteId uint) ([]SiteStatusHistory, error) {
	var status []SiteStatusHistory
	err := nsr.DB.Where("site_id = ?", siteId).Find(&status).Error
	return status, err
}

func (nsr *SiteStatusRepo) GetAllSiteStatusBySiteIdAndDate(siteId uint, firstDate, secondDate time.Time) ([]SiteStatusHistory, error) {
	var status []SiteStatusHistory

	err := nsr.DB.
		Where("site_id = ? AND created_at BETWEEN ? AND ?", siteId, firstDate, secondDate).
		Find(&status).Error
	if err != nil {
		return nil, err
	}

	return status, nil
}

func (nsr *SiteStatusRepo) GetLastSiteStatus(siteId uint) (*SiteStatusHistory, error) {
	var status SiteStatusHistory
	err := nsr.DB.Where("site_id = ?", siteId).Last(&status).Error
	if err != nil {
		return nil, err
	}
	return &status, nil
}

func (nsr *SiteStatusRepo) GetFirstSiteStatus(siteId uint) (*SiteStatusHistory, error) {
	var status SiteStatusHistory
	err := nsr.DB.Where("site_id = ?", siteId).First(&status).Error
	if err != nil {
		return nil, err
	}
	return &status, nil
}

func (nsr *SiteStatusRepo) GetAllSiteStatusByStatus(siteId uint, siteStatus string) ([]SiteStatusHistory, error) {
	var status []SiteStatusHistory
	err := nsr.DB.Where("site_id = ? AND site_status = ?", siteId, siteStatus).Find(&status).Error
	if err != nil {
		return nil, err
	}
	return status, err
}
