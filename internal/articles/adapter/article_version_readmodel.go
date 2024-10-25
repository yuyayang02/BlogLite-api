package adapter

import (
	"context"
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/qmstar0/BlogLite-api/internal/articles/application/query"
	"github.com/qmstar0/BlogLite-api/internal/articles/domain/articles"
	"github.com/qmstar0/BlogLite-api/internal/common/domain"
	"github.com/qmstar0/BlogLite-api/internal/common/e"
	"github.com/qmstar0/BlogLite-api/pkg/logging"
	"gorm.io/gorm"
	"sync"
	"time"
)

type ArticleVersion struct {
	Version     string `gorm:"primaryKey"`
	URI         string
	Source      string
	Title       string
	Description string
	Note        string
	Content     string
	CreatedAt   time.Time
}

func (a ArticleVersion) TableName() string {
	return "article_version"
}

type PostgresArticleVersionReadmodel struct {
	db   *gorm.DB
	lock *sync.Mutex
}

func NewPostgresArticleVersionReadmodel(db *gorm.DB) *PostgresArticleVersionReadmodel {
	if err := db.AutoMigrate(&ArticleVersion{}); err != nil {
		logging.Logger().Fatal("数据表初始化失败", "err", err)
	}
	return &PostgresArticleVersionReadmodel{db: db, lock: &sync.Mutex{}}
}

func (p PostgresArticleVersionReadmodel) Topics() []string {
	return []string{
		domain.EventTopic(articles.ArticleDeletedEvent{}),
		domain.EventTopic(articles.ArticleNewVersionCreatedEvent{}),
		domain.EventTopic(articles.ArticleVersionContentDeletedEvent{}),
	}
}

func (p PostgresArticleVersionReadmodel) Handle(msg *message.Message) ([]*message.Message, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	switch message.SubscribeTopicFromCtx(msg.Context()) {
	case "articles.ArticleDeletedEvent":
		return nil, p.handleArticleDeletedEvent(msg)
	case "articles.ArticleNewVersionCreatedEvent":
		return nil, p.handleArticleNewVersionCreatedEvent(msg)
	case "articles.ArticleVersionContentDeletedEvent":
		return nil, p.handleArticleVersionContentDeletedEvent(msg)
	default:
		return nil, nil
	}
}

func (p PostgresArticleVersionReadmodel) handleArticleDeletedEvent(msg *message.Message) error {
	var event articles.ArticleDeletedEvent
	if err := json.Unmarshal(msg.Payload, &event); err != nil {
		return err
	}

	err := p.db.WithContext(msg.Context()).Where("uri = ?", event.URI).Delete(&ArticleVersion{}).Error
	if err != nil {
		return err
	}
	return err
}

func (p PostgresArticleVersionReadmodel) handleArticleNewVersionCreatedEvent(msg *message.Message) error {
	var event articles.ArticleNewVersionCreatedEvent
	if err := json.Unmarshal(msg.Payload, &event); err != nil {
		return err
	}

	err := p.db.WithContext(msg.Context()).Where("version = ?", event.Version).Save(&ArticleVersion{
		URI:         event.URI,
		Version:     event.Version,
		Source:      event.Source,
		Title:       event.Title,
		Description: event.Description,
		Note:        event.Note,
		Content:     event.Content,
		CreatedAt:   event.CreatedAt,
	}).Error
	if err != nil {
		return err
	}
	return err
}

func (p PostgresArticleVersionReadmodel) handleArticleVersionContentDeletedEvent(msg *message.Message) error {
	var event articles.ArticleVersionContentDeletedEvent
	if err := json.Unmarshal(msg.Payload, &event); err != nil {
		return err
	}

	err := p.db.WithContext(msg.Context()).Delete(&ArticleVersion{
		Version: event.Version,
		URI:     event.URI,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func (p PostgresArticleVersionReadmodel) ArticleVersionList(ctx context.Context, uri string) ([]query.ArticleVersionResult, error) {
	var models = make([]ArticleVersion, 0)
	err := p.db.WithContext(ctx).Model(&ArticleVersion{}).Where("uri = ?", uri).Find(&models).Error
	if err != nil {
		return nil, e.InternalServiceError(err.Error())
	}

	var view = make([]query.ArticleVersionResult, len(models))
	for i, model := range models {
		view[i] = query.ArticleVersionResult{
			Version:   model.Version,
			Note:      model.Note,
			Title:     model.Title,
			CreatedAt: model.CreatedAt.UnixMilli(),
		}
	}
	return view, nil
}
