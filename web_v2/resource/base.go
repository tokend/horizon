package resource

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/go/signcontrol"
	"gitlab.com/tokend/horizon/db2/core"
	"gitlab.com/tokend/horizon/db2/history"
	"gitlab.com/tokend/horizon/web_v2/middleware"
	"net/http"
)

type Base struct {
	logger *logan.Entry

	CoreQ    core.QInterface
	HistoryQ history.QInterface

	Owner  string
	Signer string
}

func (b *Base) Prepare(r *http.Request) error {
	b.CoreQ = r.Context().Value(middleware.CoreQCtxKey).(core.QInterface)
	b.HistoryQ = r.Context().Value(middleware.HistoryQCtxKey).(history.QInterface)

	signer, err := signcontrol.CheckSignature(r)
	if err != nil {
		return err
	}

	b.Signer = signer

	return nil
}

func (b *Base) isSignedBy(signer string) bool {
	return b.Signer == signer
}

func (b *Base) isSignedByMaster() bool {
	master := "" // TODO: fetch master account id
	return b.isSignedBy(master)
}
