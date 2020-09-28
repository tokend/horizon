package resources

import (
	"gitlab.com/tokend/horizon/db2/history2"
	regources "gitlab.com/tokend/regources/generated"
)

func NewDeferredPaymentKey(ID int64) regources.Key {
	return regources.NewKeyInt64(ID, regources.DEFERRED_PAYMENTS)
}

func NewDeferredPayment(record history2.DeferredPayment) regources.DeferredPayment {
	return regources.DeferredPayment{
		Key: NewDeferredPaymentKey(record.ID),
		Relationships: regources.DeferredPaymentRelationships{
			Destination:   NewAccountKey(record.DestinationAccount).AsRelation(),
			Source:        NewAccountKey(record.SourceAccount).AsRelation(),
			SourceBalance: NewBalanceKey(record.SourceBalance).AsRelation(),
		},
	}
}
