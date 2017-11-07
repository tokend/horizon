package resource

import (
	"time"

	"bullioncoin.githost.io/development/go/amount"
	"bullioncoin.githost.io/development/horizon/db2/core"
	"bullioncoin.githost.io/development/horizon/db2/history"
)

func (b *BalancePublic) Populate(balance history.Balance) {
	b.BalanceID = balance.BalanceID
	b.AccountID = balance.AccountID
	b.Asset = balance.Asset
	b.ExchangeID = balance.ExchangeID
	b.ExchangeName = balance.ExchangeName
}

func (b *Balance) Populate(balance core.Balance, demurragePeriod int64) error {
	b.BalanceID = balance.BalanceID
	b.AccountID = balance.AccountID
	b.Balance = amount.String(balance.Amount + balance.Locked)
	b.Locked = amount.String(balance.Locked)
	b.StorageFee = amount.String(balance.StorageFee)
	b.StorageFeeLastCalculated = time.Unix(int64(balance.StorageFeeLastCalculated), 0).UTC()
	b.FeesPaid = amount.String(balance.FeesPaid)
	b.StorageFeeTime = time.Unix(int64(balance.StorageFeeLastCharged), 0).UTC()
	b.ExchangeID = balance.ExchangeID
	b.Asset = balance.Asset

	b.ExchangeName = balance.ExchangeName
	b.RequireReview = balance.RequireReview
	b.IncentivePerCoin = amount.String(balance.IncentivePerCoin)

	return nil
}

func (this *Balance) calcStorageFee(balance core.Balance, demurragePeriod int64) {
	// TODO rewrite!!
	/*lastCalculated := time.Unix(int64(balance.StorageFeeLastCalculated), 0).UTC()

	timeToCalcFor := time.Now().Add(time.Hour).UTC()

	var absPeriodPassed, feePeriod big.Float

	absPeriodPassed.SetFloat64(float64(int64(timeToCalcFor.Unix() - lastCalculated.Unix())))
	//feePeriod.SetFloat64(float64(balance.FeePeriod))
	// TODO use value from cfg
	feePeriod.SetFloat64(float64(demurragePeriod))


	var amountVal, feePercent, amountForWholePeriod, additionalFee big.Float

	relPeriodPassed := new(big.Float).Quo(&absPeriodPassed, &feePeriod)

	amountVal.SetInt64(balance.Amount)
	feePercent.SetInt64(balance.FeePercent)
	amountForWholePeriod.Mul(&amountVal, &feePercent).
		Quo(&amountForWholePeriod, big.NewFloat(amount.One)).
		Quo(&amountForWholePeriod, big.NewFloat(100))

	additionalFee.Mul(relPeriodPassed, &amountForWholePeriod)
	var additionalFeeVal big.Int
	additionalFee.Int(&additionalFeeVal)
	this.StorageFee = amount.String(xdr.Int64(int64(balance.StorageFee) + additionalFeeVal.Int64() - balance.FeesPaid))
	this.StorageFeeTime = timeToCalcFor*/
}

func (balance BalancePublic) PagingToken() string {
	return balance.ID
}
