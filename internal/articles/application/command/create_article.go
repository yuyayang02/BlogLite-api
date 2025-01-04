package command

import (
	"context"
	"fmt"
	"github.com/yuyayang02/BlogLite-api/internal/articles/domain/articles"
)

type CreateArticle struct {
	URI                 string
	MarkdownTextContent string
	CategoryID          string
}

type CreateArticleHandler struct {
	contentProcessingService articles.ContentProcessingService
	cateCheckService         CategoryValidityCheckService
	repo                     articles.ArticleRepository
}

func (h CreateArticleHandler) Handle(ctx context.Context, cmd CreateArticle) error {
	var err error
	uri, err := articles.NewUri(cmd.URI)
	if err != nil {
		return err
	}
	if err = h.cateCheckService.CategoryExist(ctx, cmd.CategoryID); err != nil {
		return errors.ResourceDoesNotExist.WithMessage(fmt.Sprintf("%s", err.Error()))
	}

	if find, err := h.repo.Find(ctx, uri); err != nil {
		return err
	} else if find != nil {
		return errors.InvalidActionError("资源已存在")
	}

	article, err := articles.NewArticle(uri, cmd.CategoryID, cmd.MarkdownTextContent, h.contentProcessingService)
	if err != nil {
		return err
	}

	return h.repo.Save(ctx, article)
}
