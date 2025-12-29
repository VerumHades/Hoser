package util

import (
	"context"
	"fmt"
)

/*
func WithEntity[K ~string, T any](
	ctx context.Context,
	id K,
	repo interface {
		GetByID(ctx context.Context, id K) (*T, error)
	},
	notFoundError error,
	processFunction func(ctx context.Context, entity *T) error,
) error {
	entity, err := repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if entity == nil {
		if notFoundError != nil {
			return notFoundError
		}
		return fmt.Errorf("entity %s does not exist", id)
	}
	return processFunction(ctx, entity)
}

func WithExistingEntity[K ~string](
	ctx context.Context,
	id K,
	repo interface {
		Exists(ctx context.Context, id K) (bool, error)
	},
	notFoundError error,
	processFunction func(ctx context.Context) error,
) error {
	exists, err := repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		if notFoundError != nil {
			return notFoundError
		}
		return fmt.Errorf("entity %s does not exist", id)
	}
	return processFunction(ctx)
}*/

func GetExistingEntity[ID ~string, T any](
	ctx context.Context,
	id ID,
	repo interface {
		GetByID(ctx context.Context, id ID) (*T, error)
	},
) (*T, error) {
	entity, err := repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, fmt.Errorf("entity %s does not exist", id)
	}
	return entity, nil
}

func EnsureEntityExists[ID ~string](
	ctx context.Context,
	id ID,
	notFoundError error,
	repo interface {
		Exists(ctx context.Context, id ID) (bool, error)
	},
) error {
	exists, err := repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		if notFoundError != nil {
			return notFoundError
		}
		return fmt.Errorf("entity %s does not exist", id)
	}
	return nil
}
