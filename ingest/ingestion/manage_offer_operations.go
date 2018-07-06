package ingestion

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"encoding/json"
	"github.com/guregu/null"
)

func (ingest *Ingestion) SetOfferState(offerID uint64, state string) error {
	var manageOfferOps []history.Operation
	opTypes := []xdr.OperationType{xdr.OperationTypeManageOffer}
	err := ingest.HistoryQ().Operations().ForTypes(opTypes).Select(&manageOfferOps)
	if err != nil {
		return errors.Wrap(err, "failed to load manage offer operations")
	}

	for _, op := range manageOfferOps {
		var details map[string]interface{}
		err = op.UnmarshalDetails(&details)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal manage offer operation details")
		}

		if uint64(details["offer_id"].(float64)) == offerID {
			details["offer_state"] = state
			bytes, err := json.Marshal(details)
			if err != nil {
				return errors.Wrap(err, "failed to marshal manage offer operation details")
			}
			op.DetailsString = null.StringFrom(string(bytes))
			err = ingest.HistoryQ().Operations().Update(op)
		}
	}

	return nil
}
