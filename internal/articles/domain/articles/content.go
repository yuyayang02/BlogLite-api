package articles

import (
	"fmt"
	"github.com/yuyayang02/BlogLite-api/internal/common/errors"
	"slices"
	"strings"
)

type Content struct {
	title       string
	description string
	tags        []string
	body        string
	raw         string
	hash        string
}

func (c Content) Title() string       { return c.title }
func (c Content) Description() string { return c.description }
func (c Content) Tags() []string      { return slices.Clone(c.tags) }
func (c Content) Body() string        { return c.body }
func (c Content) Raw() string         { return c.raw }
func (c Content) Hash() string        { return c.hash }

// DocumentParser 文档解析器
type DocumentParser interface {
	// Parse 从文档中解析出正文和元数据(yaml front matter)
	Parse(raw []byte) (map[string]string, string, error)
	Hash(raw []byte) string
}

// ContentFactory 内容工厂
type ContentFactory interface {
	// NewContent 从文档创建一个内容值对象
	NewContent(raw []byte) (Content, error)
}

// ContentFactoryConfig 内容工厂配置
type ContentFactoryConfig struct {
	// 内容总大小最大限制
	ContentBytesMaxSize int
	// 标签最大长度
	TitleMaxSize int
	// 简介最大长度
	DescriptionMaxSize int
	// 最大标签数量
	TagMaxNum int
	// 标签最大长度
	TagMaxSize int
}

func (c ContentFactoryConfig) Validate() error {
	if c.ContentBytesMaxSize <= 0 {
		return fmt.Errorf("`ContentBytesMaxSize` cannot be 0")
	}
	return nil
}

type contentFactory struct {
	parser DocumentParser
	config ContentFactoryConfig
}

func NewContentFactory(parser DocumentParser, config ContentFactoryConfig) (ContentFactory, error) {
	if parser == nil {
		panic("missing DocumentParser")
	}
	if err := config.Validate(); err != nil {
		return nil, err
	}
	return &contentFactory{parser, config}, nil
}

func (c *contentFactory) NewContent(raw []byte) (Content, error) {

	if len(raw) > c.config.ContentBytesMaxSize {
		return Content{}, nil
	}

	metadata, body, err := c.parser.Parse(raw)
	if err != nil {
		return Content{}, err
	}

	// 检查必需字段
	title, ok := metadata["title"]
	if !ok {
		return Content{}, errors.NewDomainError("缺少标题，请提供文档的标题（title）")
	}
	description, ok := metadata["description"]
	if !ok {
		return Content{}, errors.NewDomainError("缺少简介，请提供文档的简介（description）")
	}

	content := Content{
		raw:  string(raw),
		body: body,
		hash: c.parser.Hash(raw),
	}
	// 检查title
	content.title, err = c.CheckTitle(title)
	if err != nil {
		return Content{}, err
	}

	// 检查description
	content.description, err = c.CheckDescription(description)
	if err != nil {
		return Content{}, err
	}

	if tags, ok := metadata["tags"]; ok {
		content.tags, err = c.ParseTags(tags)
		if err != nil {
			return Content{}, err
		}
	}

	return content, nil
}

func (c *contentFactory) ParseTags(_tags string) ([]string, error) {
	// 如果标签字符串为空，返回 nil
	if _tags == "" {
		return nil, nil
	}

	// 删除空标签
	tags := slices.DeleteFunc(strings.Split(_tags, ","), func(s string) bool {
		return strings.TrimSpace(s) == ""
	})

	// 删除后的标签数量
	tagsLen := len(tags)

	// 如果删除空标签后标签数量为 0，返回 nil
	if tagsLen == 0 {
		return nil, nil
	}

	// 检查标签数量限制
	if c.config.TagMaxNum > 0 && tagsLen > c.config.TagMaxNum {
		return nil, errors.NewDomainErrorf("标签数量过多，最大允许 %d 个标签", c.config.TagMaxNum)
	}

	// 检查每个标签的长度是否超限
	filterTagMaxSize := c.config.TagMaxSize > 0
	if slices.IndexFunc(tags, func(s string) bool {
		return filterTagMaxSize && len(s) > c.config.TagMaxSize
	}) > -1 {
		return nil, errors.NewDomainErrorf("标签过长，每个标签最大长度为%d个字节", c.config.TitleMaxSize)
	}

	return tags, nil
}

func (c *contentFactory) CheckDescription(desc string) (string, error) {
	desc = strings.TrimSpace(desc)
	if c.config.DescriptionMaxSize > 0 && len(desc) > c.config.DescriptionMaxSize {
		return "", errors.NewDomainErrorf("简介过长，简介最大长度为 %d 个字节", c.config.DescriptionMaxSize)
	}
	return desc, nil
}

func (c *contentFactory) CheckTitle(title string) (string, error) {
	title = strings.TrimSpace(title)
	if c.config.TitleMaxSize > 0 && len(title) > c.config.TitleMaxSize {
		return "", errors.NewDomainErrorf("标题过长，标题最大长度为 %d 个字节", c.config.TitleMaxSize)
	}
	return title, nil
}
