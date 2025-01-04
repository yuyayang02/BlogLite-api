package service

import "context"

type MockCategoryDeletionChecker struct {
}

func (m MockCategoryDeletionChecker) CanDelete(ctx context.Context, slug string) (bool, error) {
	return true, nil
}
