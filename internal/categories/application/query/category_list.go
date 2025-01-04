package query

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/yuyayang02/BlogLite-api/internal/common/decorator"
)

type CategroyResult struct {
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CategoryListQuery struct {
}

type CategoryListResult struct {
	Count int              `json:"count"`
	Items []CategroyResult `json:"items"`
}

type (
	CategoryListReadmodel interface {
		CategoryList(context.Context) ([]CategroyResult, error)
	}

	CategoryListHandler decorator.QueryHandler[CategoryListQuery, CategoryListResult]
)

type categoryListHandler struct {
	rm CategoryListReadmodel
}

func NewCategoryListHandler(
	rm CategoryListReadmodel,
	logger *logrus.Entry,
	client decorator.MetricsClient,
) CategoryListHandler {
	return decorator.ApplyQueryDecorators[CategoryListQuery, CategoryListResult](
		&categoryListHandler{rm: rm},
		logger,
		client,
	)
}

func (h *categoryListHandler) Handle(ctx context.Context, _ CategoryListQuery) (CategoryListResult, error) {
	list, err := h.rm.CategoryList(ctx)
	if err != nil {
		return CategoryListResult{}, err
	}

	return CategoryListResult{
		Count: len(list),
		Items: list,
	}, nil
}
