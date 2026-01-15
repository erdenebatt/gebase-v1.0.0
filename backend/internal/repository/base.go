package repository

import (
	"context"

	"gorm.io/gorm"
)

type BaseRepository[T any] struct {
	DB *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{DB: db}
}

type PaginationParams struct {
	Page     int
	PageSize int
}

func (p *PaginationParams) GetOffset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 20
	}
	return (p.Page - 1) * p.PageSize
}

func (p *PaginationParams) GetLimit() int {
	if p.PageSize < 1 {
		p.PageSize = 20
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	return p.PageSize
}

type PaginatedResult[T any] struct {
	Data       []T   `json:"data"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Create(entity).Error
}

func (r *BaseRepository[T]) FindByID(ctx context.Context, id interface{}) (*T, error) {
	var entity T
	err := r.DB.WithContext(ctx).First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepository[T]) FindAll(ctx context.Context) ([]T, error) {
	var entities []T
	err := r.DB.WithContext(ctx).Find(&entities).Error
	return entities, err
}

func (r *BaseRepository[T]) FindWithPagination(ctx context.Context, params PaginationParams) (*PaginatedResult[T], error) {
	var entities []T
	var total int64

	query := r.DB.WithContext(ctx).Model(new(T))

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := query.Offset(params.GetOffset()).Limit(params.GetLimit()).Find(&entities).Error; err != nil {
		return nil, err
	}

	totalPages := int(total) / params.GetLimit()
	if int(total)%params.GetLimit() > 0 {
		totalPages++
	}

	return &PaginatedResult[T]{
		Data:       entities,
		Page:       params.Page,
		PageSize:   params.GetLimit(),
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (r *BaseRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Save(entity).Error
}

func (r *BaseRepository[T]) Delete(ctx context.Context, id interface{}) error {
	var entity T
	return r.DB.WithContext(ctx).Delete(&entity, id).Error
}

func (r *BaseRepository[T]) FindByCondition(ctx context.Context, condition map[string]interface{}) ([]T, error) {
	var entities []T
	err := r.DB.WithContext(ctx).Where(condition).Find(&entities).Error
	return entities, err
}

func (r *BaseRepository[T]) FindOneByCondition(ctx context.Context, condition map[string]interface{}) (*T, error) {
	var entity T
	err := r.DB.WithContext(ctx).Where(condition).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepository[T]) Exists(ctx context.Context, condition map[string]interface{}) (bool, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(new(T)).Where(condition).Count(&count).Error
	return count > 0, err
}
