package ctx

import (
	"net/http"

	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/go/doorman"
	"gitlab.com/tokend/horizon/corer"
	"gitlab.com/tokend/horizon/db2"
)

type ctxKey int

const (
	keyCoreRepo ctxKey = iota
	keyHistoryRepo
	keyLog
	keyDoorman
	keyCoreInfo
)

// Log - gets entry from context
func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(keyLog).(*logan.Entry)
}

// SetLog - sets log entry into ctx
func SetLog(ctx context.Context, value *logan.Entry) context.Context {
	return context.WithValue(ctx, keyLog, value)
}

// CoreRepo - returns new copy of repo with connection to core DB
func CoreRepo(r *http.Request) *db2.Repo {
	return getRepo(r, keyCoreRepo)
}

// SetCoreRepo - sets core repo which be used as source for CoreRepo
func SetCoreRepo(repo *db2.Repo) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, keyCoreRepo, repo)
	}
}

// HistoryRepo - returns new copy of repo with connection to hisotry DB
func HistoryRepo(r *http.Request) *db2.Repo {
	return getRepo(r, keyHistoryRepo)
}

// SetHistoryRepo - sets history repo which be used as source for HistoryRepo
func SetHistoryRepo(repo *db2.Repo) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, keyHistoryRepo, repo)
	}
}

func getRepo(r *http.Request, key ctxKey) *db2.Repo {
	repo := r.Context().Value(key).(*db2.Repo)
	return &db2.Repo{
		DB:  repo.DB,
		Log: Log(r),
	}
}

//Doorman - perform signature check
func Doorman(r *http.Request, constraints ...doorman.SignerConstraint) error {
	d := r.Context().Value(keyDoorman).(doorman.Doorman)
	return d.Check(r, constraints...)
}

//SetDoorman - adds doorman to context
func SetDoorman(d doorman.Doorman) func(ctx context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, keyDoorman, d)
	}
}

//SetCoreInfo - adds core info to context
func SetCoreInfo(info corer.Info) func(ctx context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, keyCoreInfo, info)
	}
}

//CoreInfo - returns core info from the context
func CoreInfo(r *http.Request) corer.Info {
	return r.Context().Value(keyCoreInfo).(corer.Info)
}
