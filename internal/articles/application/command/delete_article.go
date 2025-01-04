package command

import (
	"context"
	"github.com/yuyayang02/BlogLite-api/internal/articles/domain/articles"
)

type DeleteArticle struct {
	URI string
}

type DeleteArticleHandler struct {
	repo articles.ArticleRepository
}

func NewDeleteArticleHandler(resp articles.ArticleRepository) *DeleteArticleHandler {
	return &DeleteArticleHandler{repo: resp}
}

func (h DeleteArticleHandler) Handle(ctx context.Context, cmd DeleteArticle) error {
	uri, err := articles.NewUri(cmd.URI)
	if err != nil {
		return err
	}

	if found, err := h.repo.Find(ctx, uri); err != nil {
		return err
	} else if found == nil {
		return errors.ResourceDoesNotExist
	} else {
		found.Delete()
		return h.repo.Remove(ctx, found)
	}
}
