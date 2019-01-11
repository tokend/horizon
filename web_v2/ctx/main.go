package ctx

import (
	"net/http"

	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/horizon/db2"
)

type ctxKey int

const (
	coreRepoKey ctxKey = iota
	historyRepoKey
	signCheckSkipKey
	logKey
)

// Log - gets entry from context
func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logKey).(*logan.Entry)
}

// SetLog - sets log entry into ctx
func SetLog(ctx context.Context, value *logan.Entry) context.Context {
	return context.WithValue(ctx, logKey, value)
}

// SignCheckSkip - gets from ctx if request signature verification should be skipped
func SignCheckSkip(r *http.Request) bool {
	return r.Context().Value(signCheckSkipKey).(bool)
}

// SetSignCheckSkip - sets into context if request signature verification should be skipped
func SetSignCheckSkip(value bool) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, signCheckSkipKey, value)
	}
}

// CoreRepo - returns new copy of repo with connection to core DB
func CoreRepo(r *http.Request) *db2.Repo {
	return getRepo(r, coreRepoKey)
}

// SetCoreRepo - sets core repo which be used as source for CoreRepo
func SetCoreRepo(repo *db2.Repo) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, coreRepoKey, repo)
	}
}

// HistoryRepo - returns new copy of repo with connection to hisotry DB
func HistoryRepo(r *http.Request) *db2.Repo {
	return getRepo(r, historyRepoKey)
}

// SetHistoryRepo - sets history repo which be used as source for HistoryRepo
func SetHistoryRepo(repo *db2.Repo) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, historyRepoKey, repo)
	}
}

func getRepo(r *http.Request, key ctxKey) *db2.Repo {
	repo := r.Context().Value(key).(*db2.Repo)
	return &db2.Repo{
		DB:  repo.DB,
		Log: Log(r),
	}
}
