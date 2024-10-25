package query

import "context"

type ArticleContent struct {
	URI     string
	Version *string
}

type ArticleContentReadmodel interface {
	ArticleContent(ctx context.Context, uri string, version *string) (ArticleResult, error)
}

type ArticleContentHandler struct {
	rm ArticleContentReadmodel
}

func NewArticleContentHandler(rm ArticleContentReadmodel) *ArticleContentHandler {
	return &ArticleContentHandler{rm: rm}
}

func (h ArticleContentHandler) Handle(ctx context.Context, query ArticleContent) (ArticleResult, error) {
	return h.rm.ArticleContent(ctx, query.URI, query.Version)
}
