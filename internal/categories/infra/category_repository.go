package infra

import (
	"context"
	"github.com/yuyayang02/BlogLite-api/internal/categories/application/query"
	"github.com/yuyayang02/BlogLite-api/internal/categories/domain/categories"
	"github.com/yuyayang02/BlogLite-api/internal/common/logging"
	"gorm.io/gorm"
)

type Category struct {
	Slug string `gorm:"primaryKey"`
	Name string
}

func (c Category) TableName() string {
	return "bloglite_category"
}

type PostgresCategoryRepository struct {
	db *gorm.DB
}

func NewPostgresCategoryRepository(db *gorm.DB) *PostgresCategoryRepository {
	if err := db.AutoMigrate(&Category{}); err != nil {
		logging.Logger().Fatal("数据表初始化失败", "err", err)
	}
	return &PostgresCategoryRepository{db: db}
}

func (p PostgresCategoryRepository) Save(ctx context.Context, category *categories.Category) error {
	err := p.db.WithContext(ctx).Where("slug = ?", category.Slug()).Save(&Category{
		Slug: category.Slug(),
		Name: category.Name(),
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func (p PostgresCategoryRepository) Find(ctx context.Context, slug string) (*categories.Category, error) {
	var model Category
	result := p.db.WithContext(ctx).Where("slug = ?", slug).Limit(1).Find(&model)
	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected != 1 {
		return nil, nil
	}
	return categories.UnmarshalCategoryFromDatabase(model.Slug, model.Name), nil
}

func (p PostgresCategoryRepository) CategoryNameIsUsed(ctx context.Context, name string) (bool, error) {
	var count int64
	err := p.db.WithContext(ctx).Model(&Category{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count != 0, nil
}

func (p PostgresCategoryRepository) Remove(ctx context.Context, category *categories.Category) error {
	return p.db.WithContext(ctx).Where("slug = ?", category.Slug()).Delete(&Category{}).Error
}

func (p PostgresCategoryRepository) CategoryList(ctx context.Context) ([]query.CategroyResult, error) {
	var models = make([]Category, 0)
	err := p.db.WithContext(ctx).Find(&models).Error
	if err != nil {
		return nil, err
	}

	var views = make([]query.CategroyResult, 0, len(models))
	for _, model := range models {
		views = append(views, query.CategroyResult{
			Slug: model.Slug,
			Name: model.Name,
		})
	}

	return views, nil
}
