package ingest

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/resource/base"
)

func manageForfeitRequestToForfeitTimes(result xdr.ManageForfeitRequestResult) []base.ForfeitItem {
	items := make([]base.ForfeitItem, len(result.ForfeitRequestDetails.Items))
	for i := range result.ForfeitRequestDetails.Items {
		items[i].UnitSize = amount.String(int64(result.ForfeitRequestDetails.Items[i].Form.Unit))
		items[i].Name = string(result.ForfeitRequestDetails.Items[i].Form.Name)
		items[i].UnitsNum = int64(result.ForfeitRequestDetails.Items[i].Quantity)
	}

	return items
}
