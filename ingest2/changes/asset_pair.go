package changes

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	history "gitlab.com/tokend/horizon/db2/history2"
)

//go:generate mockery -case underscore -name assetPairStorage -inpkg -testonly
type assetPairStorage interface {
	InsertAssetPair(history.AssetPair) error
}

type assetPairHandler struct {
	storage assetPairStorage
}

func newAssetPairHandler(storage assetPairStorage) *assetPairHandler {
	return &assetPairHandler{
		storage: storage,
	}
}

func (p *assetPairHandler) Created(lc ledgerChange) error {
	assetPair := lc.LedgerChange.MustCreated().Data.MustAssetPair()
	newAssetPair := history.NewAssetPair(string(assetPair.Base), string(assetPair.Quote), int64(assetPair.CurrentPrice), lc.LedgerCloseTime)
	err := p.storage.InsertAssetPair(newAssetPair)
	if err != nil {
		return errors.Wrap(err, "failed to insert asset pair")
	}
	return nil
}
func (p *assetPairHandler) Updated(lc ledgerChange) error {
	assetPair := lc.LedgerChange.MustUpdated().Data.MustAssetPair()
	newAssetPair := history.NewAssetPair(string(assetPair.Base), string(assetPair.Quote), int64(assetPair.CurrentPrice), lc.LedgerCloseTime)
	err := p.storage.InsertAssetPair(newAssetPair)
	if err != nil {
		return errors.Wrap(err, "failed to insert asset pair")
	}
	return nil
}
