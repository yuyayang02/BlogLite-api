package command

import "context"

type CategoryDeletionChecker interface {
	CanDelete(ctx context.Context, slug string) (bool, error)
}
