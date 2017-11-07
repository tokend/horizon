package ingest

import (
	"bullioncoin.githost.io/development/go/amount"
	"bullioncoin.githost.io/development/go/xdr"
	"bullioncoin.githost.io/development/horizon/resource/base"
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
