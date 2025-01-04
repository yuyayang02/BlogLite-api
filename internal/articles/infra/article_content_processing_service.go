package infra

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuyayang02/BlogLite-api/internal/articles/domain/articles"
	"github.com/yuyayang02/BlogLite-api/internal/common/utils"
	"gopkg.in/yaml.v3"
	"strings"
)

type markdownMetadata struct {
	Title       string   `yaml:"title"`
	Tags        []string `yaml:"tags"`
	Description string   `yaml:"description"`
}

var dmp = diffmatchpatch.New()

var defaultMarkdownToHTMLTool = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,
		highlighting.NewHighlighting(
			highlighting.WithStyle("github"),
		),
		extension.CJK,
	),
	goldmark.WithParserOptions(
		parser.WithAttribute(),
	),

	goldmark.WithRendererOptions(
		html.WithXHTML(),
		html.WithHardWraps(),
	),
)

type ArticleContentProcessingServiceImpl struct {
}

func (a ArticleContentProcessingServiceImpl) GenerateDiff(oldContent, newContent articles.Content) (diff string, err error) {
	patchs := dmp.PatchMake(newContent.PatchableText(), oldContent.PatchableText())
	return dmp.PatchToText(patchs), nil
}

func (a ArticleContentProcessingServiceImpl) ApplyDiff(_content articles.Content, versions []articles.Version) (articles.Content, error) {

	var result = _content.PatchableText()
	for _, version := range versions {
		patches, err := dmp.PatchFromText(version.Diff())
		if err != nil {
			return articles.Content{}, err
		}
		result, _ = dmp.PatchApply(patches, result)
	}

	return articles.ParseFromPatchableText(result)
}

func (a ArticleContentProcessingServiceImpl) Render(content string) (string, error) {
	var buf bytes.Buffer
	err := defaultMarkdownToHTMLTool.Convert([]byte(content), &buf)
	if err != nil {
		return "", errors.New("markdown转html失败")
	}
	return buf.String(), nil
}

func (a ArticleContentProcessingServiceImpl) Parse(_content articles.MarkdownTextContent) (articles.Content, error) {
	var (
		err      error
		metadata markdownMetadata
	)

	parts := strings.SplitN(strings.TrimSpace(_content.String()), "---\n", 3)
	if len(parts) != 3 || parts[0] != "" {
		return articles.Content{}, ErrMissingHeadMetadata
	}

	// yaml解析遇到`\t`字符会报错，所以将`\t`替换为四个空格即可
	err = yaml.Unmarshal([]byte(strings.Replace(parts[1], "\t", "    ", -1)), &metadata)
	if err != nil {
		return articles.Content{}, fmt.Errorf("error parsing Markdown YAML front matter: (%w)", err)
	}

	return articles.NewContentFromService(metadata.Title, parts[2], metadata.Description, metadata.Tags)
}

func (a ArticleContentProcessingServiceImpl) Hash(content articles.Content) string {
	return utils.ShortHash(content.PatchableText())
}

var ErrMissingHeadMetadata = errors.New("missing header metadata")
