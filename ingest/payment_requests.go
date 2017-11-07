package ingest

import (
	"encoding/json"
	"time"

	"bullioncoin.githost.io/development/go/amount"
	"bullioncoin.githost.io/development/go/xdr"
	"gitlab.com/distributed_lab/tokend/horizon/db2/core"
	"gitlab.com/distributed_lab/tokend/horizon/db2/history"
	"gitlab.com/distributed_lab/tokend/horizon/resource/operations"
	sq "github.com/lann/squirrel"
)

func (ingest *Ingestion) InsertPaymentRequests(
	requestsInfo []xdr.PaymentRequestInfo,
) error {
	sql := ingest.payment_requests
	for _, requestInfo := range requestsInfo {
		details := operations.BasePayment{
			FromBalance:           requestInfo.PaymentRequest.SourceBalance.AsString(),
			ToBalance:             requestInfo.PaymentRequest.DestinationBalance.AsString(),
			From:                  requestInfo.Source.Address(),
			To:                    requestInfo.Destination.Address(),
			Amount:                amount.String(int64(requestInfo.PaymentRequest.SourceSend)),
			SourcePaymentFee:      amount.String(0),
			DestinationPaymentFee: amount.String(0),
			SourceFixedFee:        amount.String(0),
			DestinationFixedFee:   amount.String(0),
			SourcePaysForDest:     false,
		}
		djson, err := json.Marshal(details)
		if err != nil {
			return err
		}

		sql = sql.Values(
			requestInfo.PaymentRequest.PaymentId,
			requestInfo.PaymentRequest.Exchange.Address(),
			nil,
			djson,
			time.Now().UTC(),
			time.Now().UTC(),
			xdr.RequestTypeRequestTypePayment,
		)
	}

	_, err := ingest.DB.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

func (ingest *Ingestion) InsertPaymentRequest(
	ledger *core.LedgerHeader,
	paymentID uint64,
	exchange string,
	details interface{},
	accepted *bool,
	requestType xdr.RequestType,
) error {
	ledgerCloseTime := time.Unix(ledger.CloseTime, 0)
	djson, err := json.Marshal(details)
	if err != nil {
		return err
	}

	sql := ingest.payment_requests.Values(
		paymentID,
		exchange,
		accepted,
		djson,
		ledgerCloseTime,
		ledgerCloseTime,
		int(requestType),
	)

	_, err = ingest.DB.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

func (ingest *Ingestion) UpdatePaymentRequest(
	ledger *core.LedgerHeader,
	paymentID uint64,
	exchange string,
	accept bool,
) error {
	sql := sq.Update("history_payment_requests").SetMap(sq.Eq{
		"accepted":   accept,
		"updated_at": time.Unix(ledger.CloseTime, 0),
	}).Where("payment_id = ?", paymentID).Where("exchange = ?", exchange)

	_, err := ingest.DB.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

func (ingest *Ingestion) UpdatePayment(
	paymentID xdr.Uint64,
	accept bool,
	rejectReason *xdr.String256,
) error {
	state := history.SUCCESS
	if !accept {
		state = history.REJECTED
	}

	if rejectReason != nil {
		err := ingest.ingestRejectReason(string(*rejectReason), uint64(paymentID))
		if err != nil {
			return err
		}
	}
	sql := sq.Update("history_operations").SetMap(sq.Eq{
		"state": state,
	}).Where("identifier = ?", paymentID)

	_, err := ingest.DB.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}
