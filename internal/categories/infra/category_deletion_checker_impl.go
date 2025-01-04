package infra

import (
	"context"
	"gorm.io/gorm"
)

type CategoryDeletionCheckerImpl struct {
	db *gorm.DB
}

func NewGetCategoryUsedService(db *gorm.DB) *CategoryDeletionCheckerImpl {
	return &CategoryDeletionCheckerImpl{db: db}
}

func (g CategoryDeletionCheckerImpl) CanDelete(ctx context.Context, slug string) (bool, error) {
	var count int64
	err := g.db.WithContext(ctx).
		Table("article_detail").
		Where("category_id = ?", slug).
		Limit(1).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 0, nil
}
