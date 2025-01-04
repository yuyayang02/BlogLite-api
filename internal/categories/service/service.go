package service

import (
	"context"
	"github.com/yuyayang02/BlogLite-api/internal/categories/application"
	"github.com/yuyayang02/BlogLite-api/internal/categories/application/command"
	"github.com/yuyayang02/BlogLite-api/internal/categories/application/query"
	"github.com/yuyayang02/BlogLite-api/internal/categories/infra"
	"github.com/yuyayang02/BlogLite-api/internal/common/database/postgresql"
	"github.com/yuyayang02/BlogLite-api/internal/common/logging"
	"github.com/yuyayang02/BlogLite-api/internal/common/metrics"
)

func NewComponentTestApplication(ctx context.Context) *application.App {
	return newApplication(ctx, MockCategoryDeletionChecker{})
}

func NewApplication(ctx context.Context) *application.App {
	db := postgresql.GetDB()

	service := infra.NewGetCategoryUsedService(db)
	return newApplication(ctx, service)
}

func newApplication(ctx context.Context, service command.CategoryDeletionChecker) *application.App {

	db := postgresql.GetDB()

	repo := infra.NewPostgresCategoryRepository(db)

	metricsClient := metrics.NewMockMetrics()

	logger := logging.Logger()

	return &application.App{
		Command: application.Command{
			CreateCategory: command.NewCreateCategoryHandler(repo, logger, metricsClient),
			DeleteCategory: command.NewDeleteCategoryHandler(repo, service, logger, metricsClient),
		},
		Query: application.Query{
			CategoryList: query.NewCategoryListHandler(repo, logger, metricsClient),
		},
	}
}
