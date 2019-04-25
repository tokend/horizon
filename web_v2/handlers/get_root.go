package handlers

import (
	"net/http"
	"time"

	"gitlab.com/tokend/go/amount"

	"gitlab.com/distributed_lab/ape"

	"gitlab.com/tokend/horizon/web_v2/ctx"

	"gitlab.com/tokend/regources/generated"

	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/ledger"
)

//GetRoot - handles requests for `/v3` - returns state of the horizon
func GetRoot(w http.ResponseWriter, r *http.Request) {
	coreInfo := ctx.CoreInfo(r)
	currentState := ledger.CurrentState()
	currentTime := time.Now()
	response := regources.HorizonStateAttributes{
		Core:               ledgerInfoToState(currentState.Core),
		History:            ledgerInfoToState(currentState.History),
		HistoryV2:          ledgerInfoToState(currentState.History2),
		CurrentTime:        currentTime,
		CurrentTimeUnix:    currentTime.Unix(),
		EnvironmentName:    coreInfo.MasterExchangeName,
		MasterAccountId:    coreInfo.AdminAccountID,
		NetworkPassphrase:  coreInfo.NetworkPassphrase,
		Precision:          amount.One,
		TxExpirationPeriod: coreInfo.TxExpirationPeriod,
		XdrRevision:        xdr.Revision,
		CoreVersion:        coreInfo.CoreVersion,
	}

	ape.Render(w, regources.HorizonStateResponse{
		Data: regources.HorizonState{
			Key: regources.Key{
				ID:   currentTime.UTC().Format(time.RFC3339Nano),
				Type: regources.HORIZON_STATE,
			},
			Attributes: response,
		},
	})
}

func ledgerInfoToState(state ledger.State) regources.LedgerInfo {
	return regources.LedgerInfo{
		LastLedgerIncreaseTime: state.LastLedgerIncreaseTime,
		Latest:                 uint64(state.Latest),
		OldestOnStart:          uint64(state.OldestOnStart),
	}
}
