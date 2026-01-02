package util

import (
	"common/pkg/shared"
	"context"
)

func ResolveTransactionalContext(ctx context.Context, transaction shared.Transaction) context.Context {
	if transaction == nil {
		return ctx
	}
	return transaction.SessionContext(ctx)
}
