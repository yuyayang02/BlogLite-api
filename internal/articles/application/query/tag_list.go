package query

import "context"

type TagListReadmodel interface {
	TagList(ctx context.Context) ([]string, error)
}

type TagListHandler struct {
	rm TagListReadmodel
}

func NewTagListHandler(rm TagListReadmodel) *TagListHandler {
	return &TagListHandler{rm: rm}
}

func (t TagListHandler) Handle(ctx context.Context) (TagListResult, error) {
	list, err := t.rm.TagList(ctx)
	if err != nil {
		return TagListResult{}, err
	}
	return TagListResult{
		Count: len(list),
		Items: list,
	}, nil

}
