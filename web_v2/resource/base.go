package resource

import (
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/go/signcontrol"
	"gitlab.com/tokend/horizon/db2"
	"gitlab.com/tokend/horizon/db2/core"
	"gitlab.com/tokend/horizon/db2/history"
	"gitlab.com/tokend/horizon/web_v2/middleware"
	"net/http"
	"strconv"
)

type Base struct {
	Logger *logan.Entry

	CoreQ    core.QInterface
	HistoryQ history.QInterface

	Owner  string
	Signer string

	PageQuery db2.PageQueryV2

	SignCheckSkip bool
}

func (b *Base) GetUint64(r *http.Request, name string) (uint64, error) {
	res, err := strconv.ParseUint(chi.URLParam(r, name), 10, 64)

	if err != nil {
		return 0, errors.Wrap(err, "Failed to get "+name+"query param")
	}

	return res, nil
}

func (b *Base) GetPageQuery(r *http.Request) (*db2.PageQueryV2, error) {
	limit, err := b.GetUint64(r, "limit")
	if err != nil {
		return nil, err
	}

	page, err := b.GetUint64(r, "page")
	if err != nil {
		return nil, err
	}

	pageQuery, err := db2.NewPageQueryV2(page, limit)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to init page query")
	}

	return &pageQuery, nil
}

func (b *Base) Prepare(r *http.Request) error {
	coreRepo := r.Context().Value(middleware.CoreQCtxKey).(*db2.Repo)
	historyRepo := r.Context().Value(middleware.HistoryQCtxKey).(*db2.Repo)
	b.CoreQ = &core.Q{Repo: coreRepo}
	b.HistoryQ = &history.Q{Repo: historyRepo}

	signCheckSkip := r.Context().Value(middleware.SignCheckSkipCtxKey).(bool)
	b.SignCheckSkip = signCheckSkip
	if !signCheckSkip {
		signer, _ := signcontrol.CheckSignature(r)
		b.Signer = signer
	}

	b.Logger = r.Context().Value(middleware.LogCtxKey).(*logan.Entry)

	pageQuery, err := b.GetPageQuery(r)
	if err != nil {
		return errors.Wrap(err, "Failed to get page query")
	}

	b.PageQuery = *pageQuery

	return nil
}

func (b *Base) isSignedBy(signer string) bool {
	if b.SignCheckSkip {
		return true
	}

	// TODO: check also for signers
	return b.Signer == signer
}

func (b *Base) isSignedByMaster() bool {
	master := "" // TODO: fetch master account id
	return b.isSignedBy(master)
}
