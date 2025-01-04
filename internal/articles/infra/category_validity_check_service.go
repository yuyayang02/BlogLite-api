package infra

import (
	"context"
	"fmt"
	"github.com/yuyayang02/BlogLite-api/internal/categories/domain/categories"
)

type CategoryValidityCheckService struct {
	repo categories.CategoryRepository
}

func NewCategoryValidityCheckService(repo categories.CategoryRepository) *CategoryValidityCheckService {
	return &CategoryValidityCheckService{repo: repo}
}

func (c CategoryValidityCheckService) CategoryExist(ctx context.Context, categoryID string) error {
	find, err := c.repo.Find(ctx, categoryID)
	if err != nil {
		return err
	}
	if find == nil {
		return errors.ResourceDoesNotExist.WithMessage(fmt.Sprintf("%s: %s", "该分类不存在", categoryID))
	}
	return nil
}
