package command

import (
	"context"
	"github.com/yuyayang02/BlogLite-api/internal/articles/domain/articles"
	"github.com/yuyayang02/BlogLite-api/internal/common/errors"
)

type ChangeArticleCategory struct {
	URI           string
	NewCategoryID string
}

type ChangeArticleCategoryHandler struct {
	repo          articles.ArticleRepository
	categoryCheck CategoryValidityCheckService
}

func (h ChangeArticleCategoryHandler) Hanlde(ctx context.Context, cmd ChangeArticleCategory) error {
	uri, err := articles.NewUri(cmd.URI)
	if err != nil {
		return err
	}

	return h.repo.UpdateArticle(ctx, uri, func(article *articles.Article) (*articles.Article, error) {
		err := h.categoryCheck.CategoryExist(ctx, cmd.NewCategoryID)
		if err != nil {
			return nil, err
		}

		err = article.SetCategory(cmd.NewCategoryID)
		if err != nil {
			return nil, errors.Wrap(err, errors.ErrorCodeDomainError, "设置分类失败")
		}

		return article, nil
	})
}
