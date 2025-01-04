package service_test

//func init() {
//	postgresql.Init(postgresql.PostgresDSN(os.Getenv("DATABASE_DSN")))
//}

//
//
//func TestCategoryComponent(t *testing.T) {
//	ctx := context.TODO()
//	app := service.NewComponentTestApplication(ctx)
//
//	assert.NoError(t, app.Command.CreateCategory.Handle(ctx, command.CreateCategory{
//		Slug:        "new-categoryID",
//		Name:        "test-name",
//		description: "这是一段简介",
//	}))
//
//	assert.NoError(t, app.Command.ModifyCategoryDescription.Handle(ctx, command.ModifyCategoryDescription{
//		CategorySlug: "test-slug",
//		description:  "这是一段新的简介",
//	}))
//
//	assert.NoError(t, app.Command.DeleteCategory.Handle(ctx, command.CheckAndDeleteCategory{CategorySlug: "test-slug"}))
//}
