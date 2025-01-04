package command

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/yuyayang02/BlogLite-api/internal/categories/domain/categories"
	"github.com/yuyayang02/BlogLite-api/internal/common/decorator"
	"github.com/yuyayang02/BlogLite-api/internal/common/errors"
)

type CreateCategory struct {
	Slug string
	Name string
}

type CreateCategoryHandler decorator.CommandHandler[CreateCategory]

type createCategoryHandler struct {
	repo categories.CategoryRepository
}

func NewCreateCategoryHandler(
	repo categories.CategoryRepository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) CreateCategoryHandler {
	return decorator.ApplyCommandDecorators[CreateCategory](
		&createCategoryHandler{repo: repo},
		logger,
		metricsClient,
	)
}

func (h *createCategoryHandler) Handle(ctx context.Context, cmd CreateCategory) error {
	found, err := h.repo.Find(ctx, cmd.Slug)
	if err != nil {
		return errors.InternalServiceError(err)
	}

	if found != nil {
		return errors.New(errors.ErrorCodeResourceIsExists, "资源已存在")
	}

	exist, err := h.repo.CategoryNameIsUsed(ctx, cmd.Name)
	if err != nil {
		return errors.InternalServiceError(err)
	}

	if exist {
		return errors.NewDomainError("该分类命名已存在，尝试换一个名字？")
	}

	category, err := categories.NewCategory(cmd.Slug, cmd.Name)
	if err != nil {
		return errors.Wrap(err, errors.ErrorCodeDomainError, "创建分类时发生错误")
	}

	return h.repo.Save(ctx, category)
}
