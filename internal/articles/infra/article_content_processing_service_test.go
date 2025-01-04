package infra_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/yuyayang02/BlogLite-api/internal/articles/domain/articles"
	"github.com/yuyayang02/BlogLite-api/internal/articles/infra"
	"testing"
)

var mockMarkdownDoc1 = fmt.Sprintf(`
---
title: "这是标题"
description: "这是一段简介"

---
在vscode当前项目配置文件('.vscode/settings.json')中添加以下内容：

%s json
{
    "rust-analyzer.check.extraEnv": {
		"YOUR_ENV_VAR": "VALUE_XXX"
	}
}
%s
`, "```", "```")

var mockMarkdownDoc2 = fmt.Sprintf(`
---
title: "这是标题"
description: "这是一段简介"
---
在vscode当前项目配置文件('.vscode/settings.json')中添加以下内容：

%s json
{
    "rust-analyzer.check.extraEnv": {
		"YOUR_ENV_VAR": "VALUE_XXX"
	}
}
%s
## 原因

在vscode中开发Rust必然会用到'rust-analyzer'这个插件

这个插件默认启用的一个功能是：当文件保存时自动运行'cargo check'命令，这个功能还不能关掉，关掉的话其他正常开发行为都会受影响
`, "```", "```")

func TestContentProcessingService(t *testing.T) {
	impl := infra.ArticleContentProcessingServiceImpl{}

	markdownContent1, err := articles.NewMarkdownContent(mockMarkdownDoc1)
	assert.NoError(t, err)
	content1, err := impl.Parse(markdownContent1)
	assert.NoError(t, err)

	markdownContent2, err := articles.NewMarkdownContent(mockMarkdownDoc2)
	assert.NoError(t, err)
	content2, err := impl.Parse(markdownContent2)
	assert.NoError(t, err)
	//t.Logf("%+v", content1)
	//t.Logf("%+v", content2)

	content1Hash := impl.Hash(content1)
	//content2Hash := impl.Hash(content2)
	//t.Logf("%+v", content1Hash)
	//t.Logf("%+v", content2Hash)

	diff, err := impl.GenerateDiff(content1, content2)
	//t.Logf("%+v", diff)

	version := articles.NewVersionFromSerivce(content1Hash, diff)
	assert.Equal(t, version.Hash(), content1Hash)
	content, err := impl.ApplyDiff(content2, []articles.Version{version})
	assert.NoError(t, err)

	//t.Logf("%+v", content2Hash)

	assert.Equal(t, impl.Hash(content), content1Hash)
}
