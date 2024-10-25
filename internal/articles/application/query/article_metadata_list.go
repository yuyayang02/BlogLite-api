package query

import (
	"context"
	"github.com/qmstar0/BlogLite-api/internal/common/constant"
	"github.com/qmstar0/BlogLite-api/pkg/utils"
)

type ArticleMetadataList struct {
	Category *string
	Tags     []string
	Page     *int
	Limit    *int
}

type ArticleMetadataListReadmodel interface {
	ArticleMetadataList(ctx context.Context, offset, limit int, tags []string, categoryID *string) ([]ArticleMetadataResult, error)
}

type ArticleMetadataListHandler struct {
	rm ArticleMetadataListReadmodel
}

func NewArticleMetadataListHandler(rm ArticleMetadataListReadmodel) *ArticleMetadataListHandler {
	return &ArticleMetadataListHandler{rm: rm}
}

func (h ArticleMetadataListHandler) Handle(ctx context.Context, query ArticleMetadataList) (ArticleMetadataListResult, error) {
	var (
		page  = 1
		limit = constant.ArticleMetadataListDefaultLimit
	)

	if query.Page != nil && *query.Page > 1 {
		page = *query.Page
	}

	if query.Limit != nil && *query.Limit > 0 {
		limit = *query.Limit
	}

	list, err := h.rm.ArticleMetadataList(
		ctx,
		utils.Offset(page, limit),
		limit+1,
		query.Tags,
		query.Category,
	)

	if err != nil {
		return ArticleMetadataListResult{}, err
	}
	listLen := len(list)
	next := listLen > limit
	if next {
		list = list[:listLen-1]
	}
	return ArticleMetadataListResult{
		Count: len(list),
		Page:  page,
		Items: list,
		Prev:  page > 1,
		Next:  next,
	}, nil
}
