package mdtohtml

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v3"
	"strings"
)

type Metadata struct {
	Title       string `yaml:"title"`
	Note        string `yaml:"note"`
	Description string `yaml:"description"`
}

type ParsedDoc struct {
	Metadata Metadata
	Content  string
}

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

func Convert(mdcontent string) (string, error) {
	var buf bytes.Buffer
	err := defaultMarkdownToHTMLTool.Convert([]byte(mdcontent), &buf)
	if err != nil {
		return "", errors.New("markdown转html失败")
	}
	return buf.String(), nil
}

func Parse(mdcontent string) (*ParsedDoc, error) {
	var (
		err      error
		metadata Metadata
	)

	parts := strings.SplitN(mdcontent, "---", 3)

	if len(parts) == 3 {
		// yaml解析遇到`\t`字符会报错，替换为空格即可
		err = yaml.Unmarshal([]byte(parts[1]), &metadata)
		if err != nil {
			return nil, fmt.Errorf("error parsing Markdown YAML front matter: (%w)", err)
		}
		return &ParsedDoc{
			Metadata: metadata,
			Content:  parts[2],
		}, nil
	} else {
		return &ParsedDoc{
			Metadata: metadata,
			Content:  parts[0],
		}, nil
	}
}
