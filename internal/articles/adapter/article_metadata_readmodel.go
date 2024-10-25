package adapter

import (
	"context"
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/lib/pq"
	"github.com/qmstar0/BlogLite-api/internal/articles/application/query"
	"github.com/qmstar0/BlogLite-api/internal/articles/domain/articles"
	"github.com/qmstar0/BlogLite-api/internal/categories/adapter"
	"github.com/qmstar0/BlogLite-api/internal/common/domain"
	"github.com/qmstar0/BlogLite-api/internal/common/e"
	"github.com/qmstar0/BlogLite-api/pkg/logging"
	"gorm.io/gorm"
	"sync"
	"time"
)

type ArticleMetadata struct {
	URI                   string `gorm:"primaryKey"`
	CurrentVersion        string
	Visibility            bool
	FirstVersionCreatedAt time.Time
	CategoryID            string
}

func (a ArticleMetadata) TableName() string {
	return "article_metadata"
}

type articleMetadataWithVersion struct {
	ArticleMetadata
	ArticleVersion
	adapter.Category
	Tags pq.StringArray `gorm:"type:text[];not null"`
}

type PostgresArticleMetadataReadmodel struct {
	db   *gorm.DB
	lock *sync.Mutex
}

func NewPostgresArticleMetadataReadmodel(db *gorm.DB) *PostgresArticleMetadataReadmodel {
	if err := db.AutoMigrate(&ArticleMetadata{}); err != nil {
		logging.Logger().Fatal("数据表初始化失败", "err", err)
	}
	return &PostgresArticleMetadataReadmodel{db: db, lock: &sync.Mutex{}}
}

func (p PostgresArticleMetadataReadmodel) Topics() []string {
	return []string{
		domain.EventTopic(articles.ArticleDeletedEvent{}),
		domain.EventTopic(articles.ArticleContentSetSuccessfullyEvent{}),
		domain.EventTopic(articles.ArticleInitializedSuccessfullyEvent{}),
		domain.EventTopic(articles.ArticleVisibilityChangedEvent{}),
		domain.EventTopic(articles.ArticleFirstVersionCreatedEvent{}),
		domain.EventTopic(articles.ArticleCategoryChangedEvent{}),
	}
}

func (p PostgresArticleMetadataReadmodel) Handle(msg *message.Message) ([]*message.Message, error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	switch message.SubscribeTopicFromCtx(msg.Context()) {
	case "articles.ArticleDeletedEvent":
		return nil, p.handleArticleDeletedEvent(msg)
	case "articles.ArticleContentSetSuccessfullyEvent":
		return nil, p.handleArticleContentSetSuccessfullyEvent(msg)
	case "articles.ArticleInitializedSuccessfullyEvent":
		return nil, p.handleArticleInitializedSuccessfullyEvent(msg)
	case "articles.ArticleVisibilityChangedEvent":
		return nil, p.handleArticleVisibilityChangedEvent(msg)
	case "articles.ArticleFirstVersionCreatedEvent":
		return nil, p.handleArticleFirstVersionCreatedEvent(msg)
	case "articles.ArticleCategoryChangedEvent":
		return nil, p.handleArticleCategoryChangedEvent(msg)
	default:
		return nil, nil
	}
}

func (p PostgresArticleMetadataReadmodel) handleArticleCategoryChangedEvent(msg *message.Message) error {
	var event articles.ArticleCategoryChangedEvent
	if err := json.Unmarshal(msg.Payload, &event); err != nil {
		return err
	}

	err := p.db.WithContext(msg.Context()).
		Model(&ArticleMetadata{}).
		Where("uri = ?", event.URI).
		Update("category_id", event.NewCategoryID).Error
	if err != nil {
		return err
	}
	return nil
}

func (p PostgresArticleMetadataReadmodel) handleArticleFirstVersionCreatedEvent(msg *message.Message) error {
	var event articles.ArticleFirstVersionCreatedEvent
	if err := json.Unmarshal(msg.Payload, &event); err != nil {
		return err
	}
	err := p.db.WithContext(msg.Context()).
		Model(&ArticleMetadata{}).
		Where("uri = ?", event.URI).
		Update("first_version_created_at", event.CreatedAt).
		Update("current_version", event.Version).Error
	if err != nil {
		return err
	}
	return nil
}

func (p PostgresArticleMetadataReadmodel) handleArticleDeletedEvent(msg *message.Message) error {
	var event articles.ArticleDeletedEvent
	if err := json.Unmarshal(msg.Payload, &event); err != nil {
		return err
	}

	if err := p.db.WithContext(msg.Context()).Where("uri = ?", event.URI).Delete(&ArticleMetadata{}).Error; err != nil {
		return err
	}
	return nil
}

func (p PostgresArticleMetadataReadmodel) handleArticleContentSetSuccessfullyEvent(msg *message.Message) error {
	var event articles.ArticleContentSetSuccessfullyEvent
	if err := json.Unmarshal(msg.Payload, &event); err != nil {
		return err
	}

	err := p.db.WithContext(msg.Context()).
		Model(&ArticleMetadata{}).
		Where("uri = ?", event.URI).
		Update("current_version", event.Version).Error
	if err != nil {
		return err
	}
	return nil
}

func (p PostgresArticleMetadataReadmodel) handleArticleInitializedSuccessfullyEvent(msg *message.Message) error {
	var event articles.ArticleInitializedSuccessfullyEvent
	if err := json.Unmarshal(msg.Payload, &event); err != nil {
		return err
	}

	err := p.db.WithContext(msg.Context()).Where("uri = ?", event.URI).Save(&ArticleMetadata{
		URI:            event.URI,
		CurrentVersion: "",
		Visibility:     false,
		CategoryID:     event.CategoryID,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func (p PostgresArticleMetadataReadmodel) handleArticleVisibilityChangedEvent(msg *message.Message) error {
	var event articles.ArticleVisibilityChangedEvent
	if err := json.Unmarshal(msg.Payload, &event); err != nil {
		return err
	}

	err := p.db.WithContext(msg.Context()).Model(&ArticleMetadata{}).
		Where("uri = ?", event.URI).Update("visibility", event.Visibility).Error
	if err != nil {
		return err
	}
	return nil
}

func (p PostgresArticleMetadataReadmodel) joinsTable(db *gorm.DB) *gorm.DB {
	return db.Joins("LEFT JOIN article_tag ON article_metadata.uri = article_tag.article_uri").
		Joins("LEFT JOIN category ON article_metadata.category_id = category.slug").
		Group("article_metadata.uri, category.slug")
	//Joins("JOIN article_version ON article_version.uri = article_metadata.uri")
}

func (p PostgresArticleMetadataReadmodel) filterByTagsAndCategory(db *gorm.DB, tags []string, category *string) *gorm.DB {
	if len(tags) != 0 {
		db = db.Having("array_agg(article_tag.tag) @> ?", pq.Array(tags))
	}

	if category != nil && *category != "" {
		db = db.Where("category.slug = ?", *category)
	}
	return db
}

func (p PostgresArticleMetadataReadmodel) ArticleList(
	ctx context.Context,
	offset, limit int,
	tags []string,
	categoryID *string,
) ([]query.ArticleResult, error) {
	tx := p.joinsTable(
		p.db.WithContext(ctx).
			Model(&ArticleMetadata{}).
			Select(
				"article_metadata.uri",
				"article_metadata.first_version_created_at",
				"category.slug",
				"category.name",
				"article_version.title",
				"article_version.description",
				"array_remove(array_agg(article_tag.tag), null) AS tags",
			).
			Offset(offset).
			Limit(limit).
			Order("article_metadata.first_version_created_at desc").
			Joins("JOIN article_version ON article_version.uri = article_metadata.uri").
			Group("article_version.version"),
	).Where("article_version.version = article_metadata.current_version").
		Where("article_metadata.visibility = true")

	tx = p.filterByTagsAndCategory(tx, tags, categoryID)

	var articleList = make([]articleMetadataWithVersion, 0)

	err := tx.Find(&articleList).Error
	if err != nil {
		return nil, e.InternalServiceError(err.Error())
	}

	var views = make([]query.ArticleResult, len(articleList))
	for i, view := range articleList {
		views[i] = query.ArticleResult{
			Uri:         view.ArticleMetadata.URI,
			Title:       view.ArticleVersion.Title,
			Version:     view.ArticleMetadata.CurrentVersion,
			Description: view.ArticleVersion.Description,
			Note:        view.ArticleVersion.Note,
			Content:     view.ArticleVersion.Content,
			Visibility:  view.ArticleMetadata.Visibility,
			CreatedAt:   view.FirstVersionCreatedAt.UnixMilli(),
			Category: query.ArticleCategoryResult{
				Slug: view.Category.Slug,
				Name: view.Category.Name,
			},
			Tags: view.Tags,
		}
	}
	return views, nil
}

func (p PostgresArticleMetadataReadmodel) ArticleMetadataList(ctx context.Context, offset, limit int, tags []string, categoryID *string) ([]query.ArticleMetadataResult, error) {
	var models = make([]articleMetadataWithVersion, 0)
	tx := p.joinsTable(
		p.db.WithContext(ctx).
			Model(&ArticleMetadata{}).
			Select("article_metadata.*", "category.slug", "category.name", "array_remove(array_agg(article_tag.tag), null) AS tags").
			Offset(offset).
			Limit(limit).
			Order("article_metadata.first_version_created_at desc"),
	)

	tx = p.filterByTagsAndCategory(tx, tags, categoryID)

	if err := tx.Find(&models).Error; err != nil {
		return nil, e.InternalServiceError(err.Error())
	}

	var views = make([]query.ArticleMetadataResult, len(models))
	for i, model := range models {
		views[i] = query.ArticleMetadataResult{
			URI:        model.ArticleMetadata.URI,
			Version:    model.ArticleMetadata.CurrentVersion,
			Visibility: model.ArticleMetadata.Visibility,
			Category: query.ArticleCategoryResult{
				Slug: model.Category.Slug,
				Name: model.Category.Name,
			},
			FirstVersionCreatedAt: model.ArticleMetadata.FirstVersionCreatedAt.UnixMilli(),
			Tags:                  model.Tags,
		}
	}
	return views, nil
}

func (p PostgresArticleMetadataReadmodel) ArticleMetadata(ctx context.Context, uri string) (query.ArticleMetadataResult, error) {
	var model articleMetadataWithVersion
	result := p.joinsTable(
		p.db.WithContext(ctx).
			Model(&ArticleMetadata{}).
			Select("article_metadata.*", "category.slug", "category.name", "array_remove(array_agg(article_tag.tag), null) AS tags").
			Limit(1),
	).Where("article_metadata.uri = ?", uri).Find(&model)

	if result.Error != nil {
		return query.ArticleMetadataResult{}, e.InternalServiceError(result.Error.Error())
	} else if result.RowsAffected != 1 {
		return query.ArticleMetadataResult{}, e.ResourceDoesNotExist
	}

	return query.ArticleMetadataResult{
		URI:        model.ArticleMetadata.URI,
		Version:    model.ArticleMetadata.CurrentVersion,
		Visibility: model.ArticleMetadata.Visibility,
		Category: query.ArticleCategoryResult{
			Slug: model.Category.Slug,
			Name: model.Category.Name,
		},
		FirstVersionCreatedAt: model.ArticleMetadata.FirstVersionCreatedAt.UnixMilli(),
		Tags:                  model.Tags,
	}, nil
}

func (p PostgresArticleMetadataReadmodel) ArticleContent(ctx context.Context, uri string, version *string) (query.ArticleResult, error) {
	tx := p.joinsTable(p.db.WithContext(ctx).
		Model(&ArticleMetadata{}).
		Select(
			"article_metadata.uri",
			"article_metadata.first_version_created_at",
			"category.slug",
			"category.name",
			"article_version.title",
			"article_version.description",
			"article_version.content",
			"array_remove(array_agg(article_tag.tag), null) AS tags",
		).Limit(1).
		Joins("JOIN article_version ON article_version.uri = article_metadata.uri").
		Group("article_version.version"),
	).Where("article_metadata.uri = ?", uri).
		Where("article_metadata.visibility = true")

	if version != nil && *version != "" {
		tx = tx.Where("article_version.version = ?", *version)
	} else {
		tx = tx.Where("article_version.version = article_metadata.current_version")
	}

	var view articleMetadataWithVersion
	result := tx.Find(&view)

	if result.Error != nil {
		return query.ArticleResult{}, e.InternalServiceError(result.Error.Error())
	} else if result.RowsAffected != 1 {
		return query.ArticleResult{}, e.ResourceDoesNotExist
	}

	return query.ArticleResult{
		Uri:         view.ArticleMetadata.URI,
		Title:       view.ArticleVersion.Title,
		Version:     view.ArticleMetadata.CurrentVersion,
		Description: view.ArticleVersion.Description,
		Note:        view.ArticleVersion.Note,
		Content:     view.ArticleVersion.Content,
		Visibility:  view.ArticleMetadata.Visibility,
		CreatedAt:   view.ArticleMetadata.FirstVersionCreatedAt.UnixMilli(),
		Category: query.ArticleCategoryResult{
			Slug: view.Category.Slug,
			Name: view.Category.Name,
		},
		Tags: view.Tags,
	}, nil
}
