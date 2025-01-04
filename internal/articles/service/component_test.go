package service_test

//var markdownFileContent = `
//---
//title: 测试内容
//description: 这是一个测试内容
//note: init
//---
//# 文章标题
//
//这是一个测试用的 **Markdown** 文档
//`

//func init() {
//	postgresql.Init(postgresql.PostgresDSN(os.Getenv("DATABASE_DSN")))
//}
//
//func TestNewComponentTestApplication(t *testing.T) {
//	//t.Parallel()
//
//	ctx := context.TODO()
//	app := service.NewComponentTestApplication(ctx)
//
//	testInitializationArticle(t, app, ctx)
//	testAddNewVersion(t, app, ctx)
//
//	testChangeArticleVisibility(t, app, ctx)
//	testChangeArticleCategory(t, app, ctx)
//	testModifyArticleTags(t, app, ctx)
//	//testDeleteArticle(t, app, ctx)
//
//	<-time.Tick(time.Second)
//}
//
//func TestArticleDetail(t *testing.T) {
//	ctx := context.TODO()
//	app := service.NewComponentTestApplication(ctx)
//	// 使用已有版本号查询
//	version := "d1c120bc"
//	view, err := app.Query.ArticleMetadata.Handle(ctx, query.ArticleMetadata{
//		URI:     "test",
//		Hash: &version,
//	})
//	assert.NoError(t, err)
//	assert.Equal(t, view.URI, "test")
//	fmt.Printf("%+v\n", view)
//
//	view, err = app.Query.ArticleMetadata.Handle(ctx, query.ArticleMetadata{
//		URI: "test",
//	})
//	assert.NoError(t, err)
//	assert.Equal(t, view.URI, "test")
//	fmt.Printf("%+v\n", view)
//}
//
//func TestArticleList(t *testing.T) {
//	ctx := context.TODO()
//	app := service.NewComponentTestApplication(ctx)
//
//	view, err := app.Query.ArticleList.Handle(ctx, query.ArticleList{
//		Filter: nil,
//		Page:   nil,
//		Limit:  nil,
//	})
//	assert.NoError(t, err)
//	assert.Equal(t, view.Count, 1)
//
//	fmt.Printf("%+v\n", view)
//	filter := "category:new-categoryID;tags:code"
//	view, err = app.Query.ArticleList.Handle(ctx, query.ArticleList{
//		Filter: &filter,
//		Page:   nil,
//		Limit:  nil,
//	})
//	assert.NoError(t, err)
//	assert.Equal(t, view.Count, 1)
//	fmt.Printf("%+v\n", view)
//
//	page := 2
//	view, err = app.Query.ArticleList.Handle(ctx, query.ArticleList{
//		Filter: nil,
//		Page:   &page,
//		Limit:  nil,
//	})
//	assert.NoError(t, err)
//	assert.Equal(t, view.Count, 0)
//	fmt.Printf("%+v\n", view)
//
//	versionListView, err := app.Query.ArticleVersionList.Handle(ctx, query.ArticleVersionList{URI: "test"})
//	assert.NoError(t, err)
//	assert.Equal(t, versionListView.Count, 1)
//	fmt.Printf("%+v\n", versionListView)
//
//	tags, err := app.Query.TagList.Handle(ctx)
//	assert.NoError(t, err)
//	assert.Equal(t, tags.Count, 3)
//	fmt.Printf("%+v\n", tags)
//}
//
//func testDeleteArticle(t *testing.T, app *application.App, ctx context.Context) {
//	assert.NoError(t, app.Command.DeleteArticle.Handle(ctx, command.DeleteArticle{
//		URI: "test",
//	}))
//}
//
//func testRemoveVersion(t *testing.T, app *application.App, ctx context.Context, version string) {
//	assert.NoError(t, app.Command.RemoveVersion.Handle(ctx, command.RemoveVersion{
//		URI:     "test",
//		Hash: version,
//	}))
//}
//
//func testSetArticleVersion(t *testing.T, app *application.App, ctx context.Context, version string) {
//	assert.NoError(t, app.Command.SetArticleVersion.Handle(ctx, command.SetArticleVersion{
//		URI:     "test",
//		Hash: version,
//	}))
//}
//
//func testAddNewVersion(t *testing.T, app *application.App, ctx context.Context) {
//	assert.NoError(t, app.Command.AddNewVersion.Handle(ctx, command.AddNewVersion{
//		URI:    "test",
//		Source: markdownFileContent,
//	}))
//}
//
//func testModifyArticleTags(t *testing.T, app *application.App, ctx context.Context) {
//	assert.NoError(t, app.Command.ModifyArticleTags.Handle(ctx, command.ModifyArticleTags{
//		URI:  "test",
//		Tags: []string{"tag1", "game", "code"},
//	}))
//}
//
//func testChangeArticleCategory(t *testing.T, app *application.App, ctx context.Context) {
//	assert.NoError(t, app.Command.ChangeArticleCategory.Handle(ctx, command.ChangeArticleCategory{
//		URI:        "test",
//		CategoryID: "new-categoryID",
//	}))
//}
//
//func testChangeArticleVisibility(t *testing.T, app *application.App, ctx context.Context) {
//	assert.NoError(t, app.Command.ChangeArticleVisibility.Handle(ctx, command.ChangeArticleVisibility{
//		URI:        "test",
//		Visibility: true,
//	}))
//}
//
//func testInitializationArticle(t *testing.T, app *application.App, ctx context.Context) {
//	err := app.Command.InitializationArticle.Handle(ctx, command.InitializationArticle{
//		URI:        "test",
//		CategoryID: "categoryID",
//	})
//	if err != nil {
//		t.Log(err)
//	}
//}
