package repository

import (
	"context"
	"gorm.io/gorm"
	"time"
)

// Struct sera atualizada para que tenha gorm.Model para deixar mais simples depois
type SiteStatusHistory struct {
	ID           uint `gorm:"primaryKey"`
	SiteID       uint `gorm:"index"`
	Status       string
	StatusCode   int
	ResponseTime float64
	CheckedAt    time.Time
}

// Struct generica para que eu possa fazer as funções de paginação com qualquer tipo de entidade
// provavelmente ser mais gereralizada ainda depois para usar em outros repositorios
type PaginatedResult[T any] struct {
	Data       []T   `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

type SiteStatusRepo struct {
	DB *gorm.DB
}

func NewSiteStatusRepo(db *gorm.DB) *SiteStatusRepo {
	return &SiteStatusRepo{
		DB: db,
	}
}

func (r *SiteStatusRepo) Insert(ctx context.Context, siteID uint, status string, statusCode int, responseTime float64, checkedAt time.Time) error {
	return r.DB.WithContext(ctx).Create(&SiteStatusHistory{
		SiteID:       siteID,
		Status:       status,
		StatusCode:   statusCode,
		ResponseTime: responseTime,
		CheckedAt:    checkedAt,
	}).Error
}

// Funções de Filtragens
func (nsr *SiteStatusRepo) GetAllSiteStatusBySiteIdPaginated(siteId uint, page, pageSize int) ([]SiteStatusHistory, error) {
	var status []SiteStatusHistory

	offset := (page - 1) * pageSize

	err := nsr.DB.Where("site_id = ?", siteId).Order("checked_at desc").Limit(pageSize).Offset(offset).Find(&status).Error
	if err != nil {
		return nil, err
	}
	return status, nil
}

func (nsr *SiteStatusRepo) GetAllSiteStatusBySiteIdAndDate(siteId uint, firstDate, secondDate time.Time) ([]SiteStatusHistory, error) {
	var status []SiteStatusHistory

	err := nsr.DB.
		Where("site_id = ? AND checked_at BETWEEN ? AND ?", siteId, firstDate, secondDate).
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
