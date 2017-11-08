package ingest

import (
	"time"

	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/core"
	sq "github.com/lann/squirrel"
)

// Insert Emission Requests ingests the provided request data into a new row in the
// `history_emission_requests` table
func (ingest *Ingestion) InsertEmissionRequest(
	ledger *core.LedgerHeader,
	emissionRequest *xdr.CoinsEmissionRequestEntry,
) error {
	var approved *bool
	if emissionRequest.IsApproved {
		approved = &emissionRequest.IsApproved
	}

	ledgerCloseTime := time.Unix(ledger.CloseTime, 0)
	sql := ingest.emission_requests.Values(
		uint64(emissionRequest.RequestId),
		string(emissionRequest.Reference),
		emissionRequest.Issuer.Address(),
		emissionRequest.Receiver.AsString(),
		amount.String(int64(emissionRequest.Amount)),
		string(emissionRequest.Asset),
		approved,
		ledgerCloseTime,
		ledgerCloseTime,
	)

	_, err := ingest.DB.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

func (ingest *Ingestion) DeleteEmissionRequest(
	requestID uint64,
) error {
	sql := sq.Delete("history_emission_requests").Where("request_id = ?", requestID)

	_, err := ingest.DB.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

type preEmission struct {
	SerialNumber string `json:"serialNumber"`
	Amount       string `json:"amount"`
	Asset        string `json:"asset"`
}

func (ingest *Ingestion) ApproveEmissionRequestUpdate(
	ledger *core.LedgerHeader,
	requestID uint64,
) error {
	sql := sq.Update("history_emission_requests").SetMap(sq.Eq{
		"approved":   true,
		"updated_at": time.Unix(ledger.CloseTime, 0),
	}).Where("request_id = ?", requestID)

	_, err := ingest.DB.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

func (ingest *Ingestion) RejectEmissionRequestUpdate(
	ledger *core.LedgerHeader,
	requestID uint64,
	reason string,
) error {
	sql := sq.Update("history_emission_requests").SetMap(sq.Eq{
		"approved":   false,
		"reason":     reason,
		"updated_at": time.Unix(ledger.CloseTime, 0),
	}).Where("request_id = ?", requestID)

	_, err := ingest.DB.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

func (ingest *Ingestion) RejectEmissionRequestForIssuer(
	ledger *core.LedgerHeader,
	issuer string,
) error {
	sql := sq.Update("history_emission_requests").SetMap(sq.Eq{
		"approved":   false,
		"reason":     "Account was blocked",
		"updated_at": time.Unix(ledger.CloseTime, 0),
	}).Where("issuer = ?", issuer).Where("approved = ?", nil)

	_, err := ingest.DB.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}
