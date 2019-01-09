package middleware

import (
	"context"
	"gitlab.com/tokend/horizon/db2"
	"gitlab.com/tokend/horizon/db2/core"
	"gitlab.com/tokend/horizon/db2/history"
)

const (
	CoreQCtxKey int = iota
	HistoryQCtxKey
	SignCheckSkipKey
)

func CtxCoreQ(q core.QInterface) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, CoreQCtxKey, &db2.Repo{DB: q.GetRepo().DB, Ctx: ctx})
	}
}

func CtxHistoryQ(q history.QInterface) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, HistoryQCtxKey, &db2.Repo{DB: q.GetRepo().DB, Ctx: ctx})
	}
}

func SignCheckSkip (value bool) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, SignCheckSkipKey, value)
	}
}
