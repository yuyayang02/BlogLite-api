package query

import "context"

type CategroyResult struct {
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CategoryListView struct {
	Count int              `json:"count"`
	Items []CategroyResult `json:"items"`
}

type CategoryListReadmodel interface {
	CategoryList(context.Context) ([]CategroyResult, error)
}

type CategoryListHandler struct {
	rm CategoryListReadmodel
}

func NewCategoryListHandler(rm CategoryListReadmodel) *CategoryListHandler {
	return &CategoryListHandler{rm: rm}
}

func (h CategoryListHandler) Handle(ctx context.Context) (CategoryListView, error) {
	list, err := h.rm.CategoryList(ctx)
	if err != nil {
		return CategoryListView{}, err
	}

	return CategoryListView{
		Count: len(list),
		Items: list,
	}, nil
}
