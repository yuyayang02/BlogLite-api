package command

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/yuyayang02/BlogLite-api/internal/categories/domain/categories"
	"github.com/yuyayang02/BlogLite-api/internal/common/decorator"
	"github.com/yuyayang02/BlogLite-api/internal/common/errors"
)

type DeleteCategory struct {
	CategorySlug string
}

type DeleteCategoryHandler decorator.CommandHandler[DeleteCategory]

type deleteCategoryHandler struct {
	repo    categories.CategoryRepository
	checker CategoryDeletionChecker
}

func NewDeleteCategoryHandler(
	repo categories.CategoryRepository,
	checker CategoryDeletionChecker,
	logger *logrus.Entry,
	metricsclient decorator.MetricsClient,
) DeleteCategoryHandler {
	return decorator.ApplyCommandDecorators[DeleteCategory](
		&deleteCategoryHandler{repo: repo, checker: checker},
		logger,
		metricsclient,
	)
}

func (h *deleteCategoryHandler) Handle(ctx context.Context, cmd DeleteCategory) error {
	found, err := h.repo.Find(ctx, cmd.CategorySlug)
	if err != nil {
		return errors.InternalServiceError(err)
	} else if found == nil {
		return errors.New(errors.ErrorCodeResourceNotFound, "不存在的分类")
	}

	canDelete, err := h.checker.CanDelete(ctx, cmd.CategorySlug)
	if err != nil {
		return errors.InternalServiceError(err)
	}

	if !canDelete {
		return errors.NewDomainError("无法删除该分类，该分类下仍有文章")
	}

	return h.repo.Remove(ctx, found)
}
