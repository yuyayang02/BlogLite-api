package command

import (
	"context"
	"github.com/yuyayang02/BlogLite-api/internal/articles/domain/articles"
)

type UpdateArticleContent struct {
	URI                 string
	MarkdownTextContent string
}

type UpdateArticleContentHandler struct {
	repo                     articles.ArticleRepository
	contentProcessingService articles.ContentProcessingService
}

func (h UpdateArticleContentHandler) Handle(ctx context.Context, cmd UpdateArticleContent) error {
	uri, err := articles.NewUri(cmd.URI)
	if err != nil {
		return err
	}

	return h.repo.UpdateArticle(ctx, uri, func(article *articles.Article) (*articles.Article, error) {
		err := article.UpdateContent(cmd.MarkdownTextContent, h.contentProcessingService)
		if err != nil {
			return nil, err
		}
		return article, err
	})
}
