package service

import (
	"context"
	"github.com/qmstar0/shutdown"
	"github.com/yuyayang02/BlogLite-api/internal/articles/application"
	"github.com/yuyayang02/BlogLite-api/internal/articles/application/command"
	"github.com/yuyayang02/BlogLite-api/internal/articles/application/query"
	"github.com/yuyayang02/BlogLite-api/internal/articles/infra"
	categoryAdapter "github.com/yuyayang02/BlogLite-api/internal/categories/adapter"
	"github.com/yuyayang02/BlogLite-api/internal/common/database/postgresql"
)

func NewComponentTestApplication(ctx context.Context) *application.App {
	return newApplication(
		ctx,
		MockCategoryValidityCheckService{},
		infra.NewMarkdownParser(),
	)
}

func NewApplication(ctx context.Context) *application.App {
	db := postgresql.GetDB()
	categoryValidityCheckService := infrastructure.NewCategoryValidityCheckService(categoryAdapter.NewPostgresCategoryRepository(db))
	return newApplication(
		ctx,
		categoryValidityCheckService,
		infrastructure.NewMarkdownParser(),
	)
}

func newApplication(ctx context.Context, categoryService command.CategoryValidityCheckService, markdownService command.MarkdownParseService) *application.App {
	bus := infrastructure.NewBus()
	//bus.Register("test", mockHandler{})
	defer func() {
		bus.Run(ctx)
		shutdown.RegisterTasks(bus.RouterClose)
	}()

	db := postgresql.GetDB()

	repo := infrastructure.NewPostgresArticleRepository(db, bus)

	articleMetadataReadmodel := infrastructure.NewPostgresArticleMetadataReadmodel(db)
	articleTagReadmodel := infrastructure.NewPostgresArticleTagReadmodel(db)
	articleVersionReadmodel := infrastructure.NewPostgresArticleVersionReadmodel(db)

	bus.Register("article-detail-readmodel", articleMetadataReadmodel)
	bus.Register("article-tag-readmodel", articleTagReadmodel)
	bus.Register("article-version-readmodel", articleVersionReadmodel)

	return &application.App{
		Command: application.Command{
			InitializationArticle:   command.NewInitializationArticleHandler(repo, categoryService),
			RemoveVersion:           command.NewRemoveVersionHandler(repo),
			ModifyArticleTags:       command.NewModifyArticleTagsHandler(repo),
			ChangeArticleVisibility: command.NewChangeArticleVisibilityHandler(repo),
			DeleteArticle:           command.NewDeleteArticleHandler(repo),
			SetArticleVersion:       command.NewSetArticleVersionHandler(repo),
			AddNewVersion:           command.NewAddNewVersionHandler(repo, markdownService, infrastructure.NewArticleVersionDuplicationCheckService(db)),
			ChangeArticleCategory:   command.NewChangeArticleCategoryHandler(repo, categoryService),
		},
		Query: application.Query{
			TagList:             query.NewTagListHandler(articleTagReadmodel),
			ArticleContent:      query.NewArticleContentHandler(articleMetadataReadmodel),
			ArticleList:         query.NewArticleListHandler(articleMetadataReadmodel),
			ArticleVersionList:  query.NewArticleVersionListHandler(articleVersionReadmodel),
			ArticleMetadataList: query.NewArticleMetadataListHandler(articleMetadataReadmodel),
			ArticleMetadata:     query.NewArticleMetadatahandler(articleMetadataReadmodel),
		},
	}
}
