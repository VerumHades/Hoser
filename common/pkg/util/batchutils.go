package util

import (
	"common/pkg/shared"
	"context"
)

// BatchProcessorFunc defines the function signature to process a single item in a batch.
type BatchProcessorFunc[T any] func(ctx context.Context, item T) error

// NextBatchFunc defines the function signature to fetch the next batch of items.
type NextBatchFunc[T any, CursorType any] func(ctx context.Context, request shared.BatchRequest[CursorType]) (items []T, nextCursor CursorType, err error)

// GenerateInBatches returns a channel of items fetched batch by batch.
// Caller can range over the channel, and it closes automatically when done.
func GenerateInBatches[T any, CursorType any](
	ctx context.Context,
	maxBatchSize int,
	fetchNextBatch NextBatchFunc[T, CursorType],
) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		request := shared.BatchRequest[CursorType]{MaxBatchSize: maxBatchSize}
		for {
			items, nextCursor, err := fetchNextBatch(ctx, request)
			if err != nil {
				return
			}
			if len(items) == 0 {
				break
			}
			for _, item := range items {
				select {
				case out <- item:
				case <-ctx.Done():
					return
				}
			}
			request.Cursor = nextCursor
		}
	}()
	return out
}
