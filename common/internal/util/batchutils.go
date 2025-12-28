package util

import (
	"common/internal/shared"
	"context"
)

// BatchProcessorFunc defines the function signature to process a single item in a batch.
type BatchProcessorFunc[T any] func(ctx context.Context, item T) error

// NextBatchFunc defines the function signature to fetch the next batch of items.
type NextBatchFunc[T any] func(ctx context.Context, request shared.BatchRequest) (items []T, nextCursor shared.Cursor, err error)

// ProcessInBatches iteratively fetches and processes items in batches until no more items remain.
// T is the type of items being processed.
func ProcessInBatches[T any](
	ctx context.Context,
	maxBatchSize int,
	fetchNextBatch NextBatchFunc[T],
	processItem BatchProcessorFunc[T],
) error {
	request := shared.BatchRequest{
		MaxBatchSize: maxBatchSize,
	}
	for {
		items, nextCursor, err := fetchNextBatch(ctx, request)
		if err != nil {
			return err
		}
		if len(items) == 0 {
			break
		}
		for _, item := range items {
			if err := processItem(ctx, item); err != nil {
				return err
			}
		}
		request.Cursor = nextCursor
	}
	return nil
}

type ItemCheckFunction[T any, K any] func(ctx context.Context, item T) (K, error)

func LookupInBatches[T any, K any](
	ctx context.Context,
	maxBatchSize int,
	fetchNextBatch func(ctx context.Context, request shared.BatchRequest) ([]T, shared.Cursor, error),
	checkItem func(ctx context.Context, item T) (*K, error),
) (*K, error) {
	request := shared.BatchRequest{MaxBatchSize: maxBatchSize}

	for {
		items, nextCursor, err := fetchNextBatch(ctx, request)
		if err != nil {
			return nil, err
		}

		if len(items) == 0 {
			break
		}

		for _, item := range items {
			result, err := checkItem(ctx, item)
			if err != nil {
				return nil, err
			}
			if result != nil {
				return result, nil
			}
		}

		request.Cursor = nextCursor
	}

	return nil, nil
}
