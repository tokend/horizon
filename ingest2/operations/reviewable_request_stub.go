package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type reviewableRequestHandlerStub struct {

}

func (h *reviewableRequestHandlerStub) ParticipantsEffects(op xdr.ReviewRequestOp, res xdr.ReviewRequestSuccessResult,
	request xdr.ReviewableRequestEntry, source history2.ParticipantEffect, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{source}, nil
}


