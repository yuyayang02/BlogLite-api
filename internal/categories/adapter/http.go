package adapter

import (
	"github.com/gin-gonic/gin"
	"github.com/yuyayang02/BlogLite-api/internal/categories/application"
	"github.com/yuyayang02/BlogLite-api/internal/categories/application/command"
	"github.com/yuyayang02/BlogLite-api/internal/common/server/httpresponse"
	"github.com/yuyayang02/BlogLite-api/internal/common/utils"
)

type HttpServer struct {
	app *application.App
}

func NewHttpServer(app *application.App) *HttpServer {
	return &HttpServer{app: app}
}

func (h HttpServer) GetCategoryList(c *gin.Context) {
	view, err := h.app.Query.CategoryList.Handle(c.Request.Context())
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.Data(c, view)
}

func (h HttpServer) CreateCategory(c *gin.Context) {
	req, err := utils.BindJSON[CreateCategoryJSONRequestBody](c)
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	err = h.app.Command.CreateCategory.Handle(c.Request.Context(), command.CreateCategory{
		Slug:        req.Slug,
		Name:        req.Name,
		Description: req.Description,
	})

	httpresponse.ErrorOrOK(c, err)
}

func (h HttpServer) DeleteCategory(c *gin.Context, slug string) {
	err := h.app.Command.DeleteCategory.Handle(c.Request.Context(), command.CheckAndDeleteCategory{CategorySlug: slug})
	httpresponse.ErrorOrOK(c, err)
}

func (h HttpServer) ModifyCategoryDescription(c *gin.Context, slug string) {
	req, err := utils.BindJSON[ModifyCategoryDescriptionJSONRequestBody](c)
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	err = h.app.Command.ModifyCategoryDescription.Handle(c.Request.Context(), command.ModifyCategoryDescription{
		CategorySlug: slug,
		Description:  req.Descripition,
	})

	httpresponse.ErrorOrOK(c, err)
}
