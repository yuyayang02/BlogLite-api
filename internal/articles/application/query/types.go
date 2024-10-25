package query

type ArticleResult struct {
	Uri         string                `json:"uri"`
	Title       string                `json:"title"`
	Version     string                `json:"version,omitempty"`
	Description string                `json:"description"`
	Note        string                `json:"note,omitempty"`
	Content     string                `json:"content,omitempty"`
	Visibility  bool                  `json:"visibility,omitempty"`
	CreatedAt   int64                 `json:"createdAt"`
	Category    ArticleCategoryResult `json:"category"`
	Tags        []string              `json:"tags"`
}

type ArticleCategoryResult struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type ArticleListResult struct {
	Count int             `json:"count"`
	Page  int             `json:"page"`
	Items []ArticleResult `json:"items"`
	Prev  bool            `json:"prev"`
	Next  bool            `json:"next"`
}

type ArticleVersionResult struct {
	Version   string `json:"version"`
	Note      string `json:"note"`
	Title     string `json:"title"`
	CreatedAt int64  `json:"createdAt"`
}

type ArticleVersionListResult struct {
	Count int                    `json:"count"`
	Items []ArticleVersionResult `json:"items"`
}

type TagListResult struct {
	Count int      `json:"count"`
	Items []string `json:"items"`
}

type ArticleMetadataResult struct {
	URI                   string                `json:"uri"`
	Version               string                `json:"version"`
	Visibility            bool                  `json:"visibility"`
	Category              ArticleCategoryResult `json:"category"`
	FirstVersionCreatedAt int64                 `json:"firstVersionCreatedAt"`
	Tags                  []string              `json:"tags"`
}

type ArticleMetadataListResult struct {
	Count int                     `json:"count"`
	Page  int                     `json:"page"`
	Items []ArticleMetadataResult `json:"items"`
	Prev  bool                    `json:"prev"`
	Next  bool                    `json:"next"`
}
