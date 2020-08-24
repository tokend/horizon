package changes

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	history "gitlab.com/tokend/horizon/db2/history2"
)

//go:generate mockery -case underscore -name dataPairStorage -inpkg -testonly
type dataStorage interface {
	Insert(data history.Data) error
	Update(data history.Data) error
	Remove(dataID int64) error
}

type dataHandler struct {
	storage dataStorage
}

func newDataHandler(storage dataStorage) *dataHandler {
	return &dataHandler{
		storage: storage,
	}
}

func (h *dataHandler) Created(lc ledgerChange) error {
	rawData := lc.LedgerChange.MustCreated().Data.MustData()
	data := h.convertData(rawData)
	if err := h.storage.Insert(data); err != nil {
		return errors.Wrap(err, "failed to insert from created")
	}
	return nil
}

func (h *dataHandler) Updated(lc ledgerChange) error {
	rawData := lc.LedgerChange.MustUpdated().Data.MustData()
	data := h.convertData(rawData)
	if err := h.storage.Update(data); err != nil {
		return errors.Wrap(err, "failed to update from updated")
	}
	return nil
}

func (h *dataHandler) Removed(lc ledgerChange) error {
	op := lc.Operation.Body
	switch op.Type {
	case xdr.OperationTypeRemoveData:
		{
			if err := h.storage.Remove(int64(op.RemoveDataOp.DataId)); err != nil {
				return errors.Wrap(err, "failed to remove data by id")
			}
		}
	}

	return nil
}

func (h *dataHandler) convertData(raw xdr.DataEntry) history.Data {
	return history.Data{
		ID:    int64(raw.Id),
		Type:  int64(raw.Type),
		Value: []byte(raw.Value),
		Owner: raw.Owner.Address(),
	}
}
