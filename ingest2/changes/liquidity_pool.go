package changes

import (
	history "gitlab.com/tokend/horizon/db2/history2"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

type lpStorage interface {
	Insert(lp history.LiquidityPool) error
	Update(lp history.LiquidityPool) error
}

type lpHandler struct {
	storage lpStorage
}

func newLPHandler(storage lpStorage) *lpHandler {
	return &lpHandler{
		storage: storage,
	}
}

func (h *lpHandler) Created(lc ledgerChange) error {
	rawLP := lc.LedgerChange.MustCreated().Data.MustLiquidityPool()
	liquidityPool := h.convertLiquidityPool(rawLP)

	err := h.storage.Insert(*liquidityPool)
	if err != nil {
		return errors.Wrap(err, "failed to insert liquidity pool to DB", logan.F{
			"liquidity_pool": liquidityPool,
		})
	}

	return nil
}

func (h *lpHandler) Updated(lc ledgerChange) error {
	rawLP := lc.LedgerChange.MustUpdated().Data.MustLiquidityPool()
	liquidityPool := h.convertLiquidityPool(rawLP)

	err := h.storage.Update(*liquidityPool)
	if err != nil {
		return errors.Wrap(err, "failed to update liquidity pool", logan.F{
			"liquidity_pool": liquidityPool,
		})
	}

	return nil
}

func (h *lpHandler) convertLiquidityPool(raw xdr.LiquidityPoolEntry) *history.LiquidityPool {
	return &history.LiquidityPool{
		ID:              int64(raw.Id),
		Account:         raw.LiquidityPoolAccount.Address(),
		TokenAsset:      string(raw.LpTokenAssetCode),
		FirstBalanceID:  raw.FirstAssetBalance.AsString(),
		SecondBalanceID: raw.SecondAssetBalance.AsString(),
		TokensAmount:    regources.Amount(raw.LpTokensTotalCap),
		FirstReserve:    regources.Amount(raw.FirstReserve),
		SecondReserve:   regources.Amount(raw.SecondReserve),
	}
}
