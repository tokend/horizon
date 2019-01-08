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
	logger *logan.Entry

	CoreQ    core.QInterface
	HistoryQ history.QInterface

	Owner  string
	Signer string

	PageQuery db2.PageQueryV2
}

func (b *Base) GetUint64(r *http.Request, name string) (uint64, error) {
	res, err := strconv.ParseUint(chi.URLParam(r, name), 10, 64)

	if err != nil {
		return 0, errors.Wrap(err, "Failed to get " + name + "query param")
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
		return nil, errors.Wrap(err,"Failed to init page query")
	}

	return &pageQuery, nil
}

func (b *Base) Prepare(r *http.Request) error {
	b.CoreQ = r.Context().Value(middleware.CoreQCtxKey).(core.QInterface)
	b.HistoryQ = r.Context().Value(middleware.HistoryQCtxKey).(history.QInterface)

	signer, err := signcontrol.CheckSignature(r)
	if err != nil {
		return err
	}

	b.Signer = signer

	pageQuery, err := b.GetPageQuery(r)
	if err != nil {
		return err
	}

	b.PageQuery = *pageQuery

	return nil
}

func (b *Base) isSignedBy(signer string) bool {
	return b.Signer == signer
}

func (b *Base) isSignedByMaster() bool {
	master := "" // TODO: fetch master account id
	return b.isSignedBy(master)
}
