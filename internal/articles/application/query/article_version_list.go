package query

import "context"

type ArticleVersionList struct {
	Uri string
}

type ArticleVersionListReadmodel interface {
	ArticleVersionList(ctx context.Context, uri string) ([]ArticleVersionResult, error)
}

type ArticleVersionListHandler struct {
	rm ArticleVersionListReadmodel
}

func NewArticleVersionListHandler(rm ArticleVersionListReadmodel) *ArticleVersionListHandler {
	return &ArticleVersionListHandler{rm: rm}
}

func (a *ArticleVersionListHandler) Handle(ctx context.Context, query ArticleVersionList) (ArticleVersionListResult, error) {
	list, err := a.rm.ArticleVersionList(ctx, query.Uri)
	if err != nil {
		return ArticleVersionListResult{}, err
	}
	return ArticleVersionListResult{
		Count: len(list),
		Items: list,
	}, nil
}
