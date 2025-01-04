package command

import (
	"context"
)

type CategoryValidityCheckService interface {
	// CategoryExist 检查分类是否真实存在。文章不能被分配在一个不存在的分类下
	CategoryExist(ctx context.Context, categoryID string) error
}
