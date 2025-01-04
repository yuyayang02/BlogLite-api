package adapter

import (
	"github.com/gin-gonic/gin"
	"github.com/yuyayang02/BlogLite-api/internal/articles/application"
	"github.com/yuyayang02/BlogLite-api/internal/articles/application/command"
	"github.com/yuyayang02/BlogLite-api/internal/articles/application/query"
	"github.com/yuyayang02/BlogLite-api/internal/common/auth"
	"github.com/yuyayang02/BlogLite-api/internal/common/server/httpresponse"
	"github.com/yuyayang02/BlogLite-api/internal/common/utils"
	"io"
)

type HttpServer struct {
	app *application.App
}

func NewHttpServer(app *application.App) *HttpServer {
	return &HttpServer{app: app}
}

func (h HttpServer) GetArticleList(c *gin.Context, params GetArticleListParams) {
	view, err := h.app.Query.ArticleList.Handle(c.Request.Context(), query.ArticleList{
		Category: params.Category,
		Tags:     parseTagsFromTagParams(params.Tags),
		Page:     params.Page,
		Limit:    params.Limit,
	})

	httpresponse.ErrorOrData(c, err, view)
}

func (h HttpServer) InitializationArticle(c *gin.Context) {
	if err := auth.FilterAuthWithUserType(c.Request.Context(), auth.ADMIN); err != nil {
		httpresponse.Error(c, err)
		return
	}

	req, err := utils.BindJSON[InitializationArticleJSONRequestBody](c)
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	err = h.app.Command.InitializationArticle.Handle(c.Request.Context(), command.InitializationArticle{
		URI:        req.Uri,
		CategoryID: req.Category,
	})

	httpresponse.ErrorOrOK(c, err)
}

func (h HttpServer) DeleteArticle(c *gin.Context, uri string) {
	if err := auth.FilterAuthWithUserType(c.Request.Context(), auth.ADMIN); err != nil {
		httpresponse.Error(c, err)
		return
	}

	err := h.app.Command.DeleteArticle.Handle(c.Request.Context(), command.DeleteArticle{URI: uri})

	httpresponse.ErrorOrOK(c, err)
}

func (h HttpServer) SetArticleVersion(c *gin.Context, uri string) {
	if err := auth.FilterAuthWithUserType(c.Request.Context(), auth.ADMIN); err != nil {
		httpresponse.Error(c, err)
		return
	}

	req, err := utils.BindJSON[SetArticleVersionJSONRequestBody](c)
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	err = h.app.Command.SetArticleVersion.Handle(c.Request.Context(), command.SetArticleVersion{
		Uri:     uri,
		Version: req.Version,
	})

	httpresponse.ErrorOrOK(c, err)
}

func (h HttpServer) ChangeArticleCategory(c *gin.Context, uri string) {
	if err := auth.FilterAuthWithUserType(c.Request.Context(), auth.ADMIN); err != nil {
		httpresponse.Error(c, err)
		return
	}

	req, err := utils.BindJSON[ChangeArticleCategoryJSONRequestBody](c)
	if err != nil {
		httpresponse.Error(c, err)
		return
	}
	err = h.app.Command.ChangeArticleCategory.Handle(c.Request.Context(), command.ChangeArticleCategory{
		Uri:        uri,
		CategoryID: req.Category,
	})

	httpresponse.ErrorOrOK(c, err)
}

func (h HttpServer) ModifyArticleTags(c *gin.Context, uri string) {
	if err := auth.FilterAuthWithUserType(c.Request.Context(), auth.ADMIN); err != nil {
		httpresponse.Error(c, err)
		return
	}

	req, err := utils.BindJSON[ModifyArticleTagsJSONRequestBody](c)
	if err != nil {
		httpresponse.Error(c, err)
		return
	}
	err = h.app.Command.ModifyArticleTags.Handle(c.Request.Context(), command.ModifyArticleTags{
		Uri:  uri,
		Tags: req.Tags,
	})

	httpresponse.ErrorOrOK(c, err)
}

func (h HttpServer) CreateNewArticleVersion(c *gin.Context, uri string) {
	if err := auth.FilterAuthWithUserType(c.Request.Context(), auth.ADMIN); err != nil {
		httpresponse.Error(c, err)
		return
	}

	formFile, err := c.FormFile("content")
	if err != nil {
		httpresponse.Error(c, err)
		return
	}
	file, err := formFile.Open()
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	contentBytes, err := io.ReadAll(file)
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	err = h.app.Command.AddNewVersion.Handle(c.Request.Context(), command.AddNewVersion{
		Uri:    uri,
		Source: string(contentBytes),
	})

	httpresponse.ErrorOrOK(c, err)
}

func (h HttpServer) RemoveArticleVersion(c *gin.Context, uri string, version string) {
	if err := auth.FilterAuthWithUserType(c.Request.Context(), auth.ADMIN); err != nil {
		httpresponse.Error(c, err)
		return
	}

	err := h.app.Command.RemoveVersion.Handle(c.Request.Context(), command.RemoveVersion{
		Uri:     uri,
		Version: version,
	})

	httpresponse.ErrorOrOK(c, err)
}

func (h HttpServer) ChangeArticleVisibility(c *gin.Context, uri string) {
	if err := auth.FilterAuthWithUserType(c.Request.Context(), auth.ADMIN); err != nil {
		httpresponse.Error(c, err)
		return
	}

	req, err := utils.BindJSON[ChangeArticleVisibilityJSONRequestBody](c)
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	err = h.app.Command.ChangeArticleVisibility.Handle(c.Request.Context(), command.ChangeArticleVisibility{
		Uri:        uri,
		Visibility: req.Visibility,
	})

	httpresponse.ErrorOrOK(c, err)
}

func (h HttpServer) GetAllTags(c *gin.Context) {
	view, err := h.app.Query.TagList.Handle(c)

	httpresponse.ErrorOrData(c, err, view)
}

func (h HttpServer) GetArticleContent(c *gin.Context, uri string) {
	view, err := h.app.Query.ArticleContent.Handle(c.Request.Context(), query.ArticleContent{
		URI:     uri,
		Version: nil,
	})

	httpresponse.ErrorOrData(c, err, view)
}

func (h HttpServer) GetAuthorArticleList(c *gin.Context, params GetAuthorArticleListParams) {
	if err := auth.FilterAuthWithUserType(c.Request.Context(), auth.ADMIN); err != nil {
		httpresponse.Error(c, err)
		return
	}

	view, err := h.app.Query.ArticleMetadataList.Handle(c.Request.Context(), query.ArticleMetadataList{
		Category: params.Category,
		Tags:     parseTagsFromTagParams(params.Tags),
		Page:     params.Page,
		Limit:    params.Limit,
	})

	httpresponse.ErrorOrData(c, err, view)
}

func (h HttpServer) GetArticleMetadata(c *gin.Context, uri string) {
	if err := auth.FilterAuthWithUserType(c.Request.Context(), auth.ADMIN); err != nil {
		httpresponse.Error(c, err)
		return
	}

	view, err := h.app.Query.ArticleMetadata.Handle(c.Request.Context(), query.ArticleMetadata{URI: uri})

	httpresponse.ErrorOrData(c, err, view)
}

func (h HttpServer) GeArticleAllVersion(c *gin.Context, uri string) {
	if err := auth.FilterAuthWithUserType(c.Request.Context(), auth.ADMIN); err != nil {
		httpresponse.Error(c, err)
		return
	}

	view, err := h.app.Query.ArticleVersionList.Handle(c.Request.Context(), query.ArticleVersionList{
		Uri: uri,
	})

	httpresponse.ErrorOrData(c, err, view)
}

func (h HttpServer) GetArticleContentByVersion(c *gin.Context, uri string, version string) {
	if err := auth.FilterAuthWithUserType(c.Request.Context(), auth.ADMIN); err != nil {
		httpresponse.Error(c, err)
		return
	}

	view, err := h.app.Query.ArticleContent.Handle(c.Request.Context(), query.ArticleContent{
		URI:     uri,
		Version: &version,
	})

	httpresponse.ErrorOrData(c, err, view)
}
