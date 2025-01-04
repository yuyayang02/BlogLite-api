package command

import (
	"context"
	"github.com/yuyayang02/BlogLite-api/internal/articles/domain/articles"
)

type RevertArticleContent struct {
	URI        string
	TargetHash string
}

type RevertArticleContentHandler struct {
	repo                     articles.ArticleRepository
	contentProcessingService articles.ContentProcessingService
}

func (h RevertArticleContentHandler) Handle(ctx context.Context, cmd RevertArticleContent) error {
	uri, err := articles.NewUri(cmd.URI)
	if err != nil {
		return err
	}

	return h.repo.UpdateArticle(ctx, uri, func(article *articles.Article) (*articles.Article, error) {
		if err := article.RevertToVersionByHash(cmd.TargetHash, h.contentProcessingService); err != nil {
			return nil, err
		}

		return article, nil
	})
}
