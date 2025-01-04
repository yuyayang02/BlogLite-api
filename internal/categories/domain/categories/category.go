package categories

import "github.com/yuyayang02/BlogLite-api/internal/common/errors"

type Category struct {
	slug string
	name string
}

const (
	CategoryNameMaxSize = 30
	CategorySlugMaxSize = 20
)

func NewCategory(slug, name string) (*Category, error) {

	if len(name) > CategoryNameMaxSize {
		return nil, errors.NewDomainErrorf("分类名称过长，最大长度为 %d 字节", CategoryNameMaxSize)
	}

	if len(slug) > CategorySlugMaxSize {
		return nil, errors.NewDomainErrorf("分类标识符过长，最大长度为 %d 字节", CategorySlugMaxSize)
	}

	return &Category{
		slug: slug,
		name: name,
	}, nil
}

func (c *Category) Slug() string { return c.slug }
func (c *Category) Name() string { return c.name }

func UnmarshalCategoryFromDatabase(
	slug, name string,
) *Category {
	return &Category{
		slug: slug,
		name: name,
	}
}
