package articles

import (
	"github.com/yuyayang02/BlogLite-api/internal/common/errors"
	"slices"
	"time"
)

type Article struct {
	events []any

	uri URI

	versionList VersionHistory

	currentVersion Version

	categoryID string

	state State
}

func NewArticle(uri URI, categoryID string, content Content) *Article {
	a := &Article{
		events:         make([]any, 0, 1),
		uri:            uri,
		versionList:    newVersionHistory(),
		currentVersion: Version{hash: content.Hash()},
		categoryID:     categoryID,
		state:          PrivateState,
	}

	a.addEvent(ArticleCreatedEvent{
		URI:            uri.String(),
		CurrentVersion: content.Hash(),
		CategoryID:     categoryID,
		Tags:           content.Tags(),
		State:          PrivateState.ToInt(),
		Title:          content.Title(),
		Body:           content.Body(),
		Description:    content.Description(),
		CreatedAt:      time.Now(),
	})

	return a
}

func (a *Article) addEvent(event any) {
	a.events = append(a.events, event)
}

func (a *Article) UpdateContent(content Content) {

	previousVersion := a.currentVersion

	a.versionList = a.versionList.Archive(previousVersion)
	a.currentVersion = Version{hash: content.Hash()}

	a.addEvent(ArticleContentUpdatedEvent{
		URI:             a.uri.String(),
		PreviousVersion: previousVersion.Hash(),
		CurrentVersion:  a.currentVersion.Hash(),
		Tags:            content.Tags(),
		Title:           content.Title(),
		Body:            content.Body(),
		Description:     content.Description(),
		Updated:         time.Now(),
	})
}

func (a *Article) RevertContentByVersion(versionHash string) error {
	if a.currentVersion.Hash() == versionHash {
		return errors.NewDomainErrorf("恢复版本不能是当前版本：%s", versionHash)
	}

	history, err := a.versionList.RevertTo(versionHash)
	if err != nil {
		return err
	}

	revertedVersions := a.versionList.Diff(history)

	a.versionList = history
	a.currentVersion = Version{hash: versionHash}

	a.addEvent(ArticleContentRevertedEvent{
		URI:             a.uri.String(),
		RevertedVersion: revertedVersions,
		CurrentVersion:  a.currentVersion.Hash(),
		Updated:         time.Now(),
	})

	return nil
}
func (a *Article) Delete() { a.addEvent(ArticleDeletedEvent{URI: a.uri.String()}) }

func (a *Article) SetCategory(categoryID string) error {
	if a.categoryID == categoryID {
		return errors.NewDomainError("修改分类不能与当前分类相同")
	}

	a.addEvent(ArticleCategoryChangedEvent{
		URI:           a.uri.String(),
		OldCategoryID: a.categoryID,
		NewCategoryID: categoryID,
	})

	a.categoryID = categoryID
	return nil
}

func (a *Article) SetState(state State) error {
	if a.state == state {
		return errors.NewDomainError("无效的状态修改")
	}

	a.state = state

	a.addEvent(ArticleStateChangedEvent{
		URI:   a.uri.String(),
		State: a.state.ToInt(),
	})

	return nil
}

func GetArticleDomainEvents(article *Article) []any {
	return slices.Clone(article.events)
}
