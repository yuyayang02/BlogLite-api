package query

import (
	"context"
	"github.com/yuyayang02/BlogLite-api/internal/common/constant"
	"github.com/yuyayang02/BlogLite-api/internal/common/utils"
)

type ArticleList struct {
	Category *string
	Tags     []string
	Page     *int
	Limit    *int
}

type ArticleListReadmodel interface {
	// ArticleList 返回总文章数，当前页内容，错误
	ArticleList(ctx context.Context, offset, limit int, tags []string, categoryID *string) (int, []ArticleResult, error)
}

type ArticleListHandler struct {
	rm ArticleListReadmodel
}

func NewArticleListHandler(rm ArticleListReadmodel) *ArticleListHandler {
	return &ArticleListHandler{rm: rm}
}

func (a *ArticleListHandler) Handle(ctx context.Context, query ArticleList) (ArticleListResult, error) {
	var (
		page  = 1
		limit = constant.ArticleListDefaultLimit
	)

	if query.Page != nil && *query.Page > 1 {
		page = *query.Page
	}

	if query.Limit != nil && *query.Limit > 0 {
		limit = *query.Limit
	}

	total, list, err := a.rm.ArticleList(
		ctx,
		utils.Offset(page, limit),
		limit,
		query.Tags,
		query.Category,
	)
	if err != nil {
		return ArticleListResult{}, err
	}

	return ArticleListResult{
		Total: total,
		Page:  page,
		Items: list,
	}, nil
}
