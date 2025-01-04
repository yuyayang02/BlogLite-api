package articles

import "time"

type ArticleCreatedEvent struct {
	URI string

	CurrentVersion string

	CategoryID string
	Tags       []string

	State int

	Title       string
	Body        string
	Description string
	CreatedAt   time.Time
}

type ArticleDeletedEvent struct {
	URI string
}

type ArticleContentUpdatedEvent struct {
	URI string

	PreviousVersion string
	CurrentVersion  string

	Tags []string

	Title       string
	Body        string
	Description string
	Updated     time.Time
}

type ArticleCategoryChangedEvent struct {
	URI           string
	OldCategoryID string
	NewCategoryID string
}

type ArticleContentRevertedEvent struct {
	URI string

	RevertedVersion []string
	CurrentVersion  string

	Updated time.Time
}

type ArticleStateChangedEvent struct {
	URI   string
	State int
}
