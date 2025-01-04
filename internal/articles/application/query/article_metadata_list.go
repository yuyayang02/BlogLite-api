package query

import (
	"context"
	"github.com/yuyayang02/BlogLite-api/internal/common/constant"
	"github.com/yuyayang02/BlogLite-api/internal/common/utils"
)

type ArticleMetadataList struct {
	Category *string
	Tags     []string
	Page     *int
	Limit    *int
}

type ArticleMetadataListReadmodel interface {
	// ArticleMetadataList 返回总文章数，当前页数据，错误
	ArticleMetadataList(ctx context.Context, offset, limit int, tags []string, categoryID *string) (int64, []ArticleMetadataResult, error)
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

	total, list, err := h.rm.ArticleMetadataList(
		ctx,
		utils.Offset(page, limit),
		limit,
		query.Tags,
		query.Category,
	)
	if err != nil {
		return ArticleMetadataListResult{}, err
	}

	return ArticleMetadataListResult{
		Total: int(total),
		Page:  page,
		Items: list,
	}, nil
}
