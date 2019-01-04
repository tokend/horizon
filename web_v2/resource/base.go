package resource

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/horizon/db2/core"
	"gitlab.com/tokend/horizon/db2/history"
	"gitlab.com/tokend/horizon/web_v2/middleware"
	"net/http"
)

type Base struct {
	logger *logan.Entry

	W http.ResponseWriter
	R *http.Request

	coreQ    core.QInterface
	historyQ history.QInterface

	Owner  string
	Signer string
}

func (b *Base) CoreQ() core.QInterface {
	if b.coreQ != nil {
		return b.coreQ
	}

	b.coreQ = b.R.Context().Value(middleware.CoreQCtxKey).(core.QInterface)

	return b.coreQ
}

func (b *Base) HistoryQ() history.QInterface {
	if b.coreQ != nil {
		return b.historyQ
	}

	b.coreQ = b.R.Context().Value(middleware.HistoryQCtxKey).(core.QInterface)

	return b.historyQ
}

func (b *Base) isSignedByOwner() bool {
	return b.Signer == b.Owner
}

func (b *Base) isSignedByAdmin() bool {
	// TODO: get master signers with Q, check if our signer is master
	return false
}
