package command

import (
	"context"
	"github.com/yuyayang02/BlogLite-api/internal/articles/domain/articles"
)

type ChangeArticleState struct {
	URI    string
	Public bool
}

type ChangeArticleStateHandler struct {
	repo articles.ArticleRepository
}

func (h ChangeArticleStateHandler) Handle(ctx context.Context, cmd ChangeArticleState) error {
	uri, err := articles.NewUri(cmd.URI)
	if err != nil {
		return err
	}

	return h.repo.UpdateArticle(ctx, uri, func(article *articles.Article) (*articles.Article, error) {
		if cmd.Public {
			article.SetPublic()
		} else {
			article.SetPrivate()
		}
		return article, nil
	})
}
