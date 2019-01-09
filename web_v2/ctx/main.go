package ctx

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/horizon/db2"
	"gitlab.com/tokend/horizon/db2/core"
	"gitlab.com/tokend/horizon/db2/history"
	"gitlab.com/tokend/horizon/web_v2/middleware"
	"net/http"
)

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(middleware.LogCtxKey).(*logan.Entry)
}

func SignCheckSkip (r *http.Request) bool {
	return r.Context().Value(middleware.SignCheckSkipCtxKey).(bool)
}

func CoreQ (r *http.Request) *core.Q {
	repo := r.Context().Value(middleware.CoreQCtxKey).(*db2.Repo)
	return &core.Q{
		Repo: repo,
	}
}

func HistoryQ (r *http.Request) *history.Q {
	repo := r.Context().Value(middleware.HistoryQCtxKey).(*db2.Repo)
	return &history.Q{
		Repo: repo,
	}
}
