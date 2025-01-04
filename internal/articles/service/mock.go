package service

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/yuyayang02/BlogLite-api/internal/articles/domain/articles"
)

type MockCategoryValidityCheckService struct {
}

func (m MockCategoryValidityCheckService) CategoryExist(ctx context.Context, categoryID string) error {
	return nil
}

type mockHandler struct{}

func (h mockHandler) Topics() []string {
	return []string{
		domain.EventTopic(articles.ArticleContentSetSuccessfullyEvent{}),
		domain.EventTopic(articles.ArticleTagsModifiedEvent{}),
		domain.EventTopic(articles.ArticleInitializedSuccessfullyEvent{}),
		domain.EventTopic(articles.ArticleCategoryChangedEvent{}),
		domain.EventTopic(articles.ArticleDeletedEvent{}),
		domain.EventTopic(articles.ArticleVersionContentDeletedEvent{}),
		domain.EventTopic(articles.ArticleVisibilityChangedEvent{}),
		domain.EventTopic(articles.ArticleNewVersionCreatedEvent{}),
		domain.EventTopic(articles.ArticleFirstVersionCreatedEvent{}),
	}
}

func (h mockHandler) Handle(msg *message.Message) ([]*message.Message, error) {
	fmt.Println("mockHandler", msg.UUID, string(msg.Payload), message.PublishTopicFromCtx(msg.Context()))
	return nil, nil
}
