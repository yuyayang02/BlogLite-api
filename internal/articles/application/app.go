package application

import (
	"github.com/yuyayang02/BlogLite-api/internal/articles/application/command"
	"github.com/yuyayang02/BlogLite-api/internal/articles/application/query"
)

type App struct {
	Command Command
	Query   Query
}

type Command struct {
	InitializationArticle   *command.InitializationArticleHandler
	RemoveVersion           *command.RemoveVersionHandler
	ModifyArticleTags       *command.ModifyArticleTagsHandler
	ChangeArticleVisibility *command.ChangeArticleVisibilityHandler
	DeleteArticle           *command.DeleteArticleHandler
	SetArticleVersion       *command.SetArticleVersionHandler
	AddNewVersion           *command.AddNewVersionHandler
	ChangeArticleCategory   *command.ChangeArticleCategoryHandler
}

type Query struct {
	TagList             *query.TagListHandler
	ArticleContent      *query.ArticleContentHandler
	ArticleList         *query.ArticleListHandler
	ArticleVersionList  *query.ArticleVersionListHandler
	ArticleMetadataList *query.ArticleMetadataListHandler
	ArticleMetadata     *query.ArticleMetadatahandler
}
