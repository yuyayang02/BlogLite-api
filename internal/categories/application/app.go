package application

import (
	"github.com/yuyayang02/BlogLite-api/internal/categories/application/command"
	"github.com/yuyayang02/BlogLite-api/internal/categories/application/query"
)

type App struct {
	Command Command
	Query   Query
}

type Command struct {
	CreateCategory command.CreateCategoryHandler
	DeleteCategory command.DeleteCategoryHandler
}

type Query struct {
	CategoryList query.CategoryListHandler
}
