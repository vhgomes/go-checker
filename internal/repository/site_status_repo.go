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

// Esse aqui é um outro helper para que não repita codigo toda hora, eu acredito que seja generico e capaz de poder usar
// em outras estruturas também.
func PaginateQuery[T any](db *gorm.DB, page, pageSize int, dest *[]T) (PaginatedResult[T], error) {
	var total int64
	db.Count(&total)

	offset := (page - 1) * pageSize
	err := db.Offset(offset).Limit(pageSize).Find(dest).Error
	if err != nil {
		return PaginatedResult[T]{}, err
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return PaginatedResult[T]{
		Data:       *dest,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// ===================================================================================================//
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
func (nsr *SiteStatusRepo) GetAllSiteStatusBySiteIdPaginated(userId, siteId uint, page, pageSize int) (PaginatedResult[SiteStatusHistory], error) {
	var status []SiteStatusHistory
	query := nsr.DB.Model(&SiteStatusHistory{}).
		Joins("JOIN sites ON sites.id = site_status_histories.site_id").
		Where("site_status_histories.site_id = ? AND sites.user_id = ?", siteId, userId).
		Order("site_status_histories.checked_at DESC")

	return PaginateQuery(query, page, pageSize, &status)
}

func (nsr *SiteStatusRepo) GetAllSiteStatusBySiteIdAndDatePaginated(userId, siteId uint, firstDate, secondDate time.Time, page, pageSize int) (PaginatedResult[SiteStatusHistory], error) {
	var status []SiteStatusHistory
	query := nsr.DB.Model(&SiteStatusHistory{}).
		Joins("JOIN sites ON sites.id = site_status_histories.site_id").
		Where("site_status_histories.site_id = ? AND sites.user_id = ? AND site_status_histories.checked_at BETWEEN ? AND ?",
			siteId, userId, firstDate, secondDate).
		Order("site_status_histories.checked_at DESC")

	return PaginateQuery(query, page, pageSize, &status)
}

func (nsr *SiteStatusRepo) GetAllSiteStatusByStatusPaginated(userId, siteId uint, siteStatus string, page, pageSize int) (PaginatedResult[SiteStatusHistory], error) {
	var status []SiteStatusHistory
	query := nsr.DB.Model(&SiteStatusHistory{}).
		Joins("JOIN sites ON sites.id = site_status_histories.site_id").
		Where("site_status_histories.site_id = ? AND sites.user_id = ? AND site_status_histories.status = ?",
			siteId, userId, siteStatus).
		Order("site_status_histories.checked_at DESC")

	return PaginateQuery(query, page, pageSize, &status)
}

// Essas duas funções eu posso generalizar para só uma porém analisando superficialmente eu acho melhor ficar separado
func (nsr *SiteStatusRepo) GetLastSiteStatus(userId, siteId uint) (*SiteStatusHistory, error) {
	var status SiteStatusHistory
	err := nsr.DB.Model(&SiteStatusHistory{}).
		Joins("JOIN sites ON sites.id = site_status_histories.site_id").
		Where("site_status_histories.site_id = ? AND sites.user_id = ?", siteId, userId).
		Order("site_status_histories.checked_at DESC").
		First(&status).Error

	if err != nil {
		return nil, err
	}
	return &status, nil
}

func (nsr *SiteStatusRepo) GetFirstSiteStatus(userId, siteId uint) (*SiteStatusHistory, error) {
	var status SiteStatusHistory
	err := nsr.DB.Model(&SiteStatusHistory{}).
		Joins("JOIN sites ON sites.id = site_status_histories.site_id").
		Where("site_status_histories.site_id = ? AND sites.user_id = ?", siteId, userId).
		Order("site_status_histories.checked_at ASC").
		First(&status).Error

	if err != nil {
		return nil, err
	}
	return &status, nil
}
