package changes

import (
	"a5ac0aaf.ngrok.io/logan/v3"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

// Consumer - consumes ledger changes

type creatable interface {
	//Handles Created Ledger entry change
	Created(change LedgerChange) error
}
type updatable interface {
	//Handles Updated Ledger entry change
	Updated(change LedgerChange) error
}
type deletable interface {
	//Handles Removed Ledger entry change
	Deleted(change LedgerChange) error
}

type Consumer struct {
	Delete map[xdr.LedgerEntryType]deletable
	Update map[xdr.LedgerEntryType]updatable
	Create map[xdr.LedgerEntryType]creatable
}

func (c *Consumer) Consume(lc LedgerChange) error {
	switch lc.LedgerChange.Type {
	case xdr.LedgerEntryChangeTypeCreated:
		return c.Created(lc)
	case xdr.LedgerEntryChangeTypeUpdated:
		return c.Updated(lc)
	case xdr.LedgerEntryChangeTypeRemoved:
		return c.Deleted(lc)
	case xdr.LedgerEntryChangeTypeState:
		return nil
	default:
		return errors.From(errors.New("Unrecognized ledger entry change type"), logan.F{
			"change_type": lc.LedgerChange.Type,
		})
	}
}

func (c *Consumer) Created(lc LedgerChange) error {
	handler, ok := c.Create[lc.LedgerChange.Created.Data.Type]
	if !ok {
		return errors.From(errors.New("No handler provided for required entry type"), logan.F{
			"entry_type": lc.LedgerChange.Created.Data.Type,
		})
	}

	return handler.Created(lc)
}

func (c *Consumer) Updated(lc LedgerChange) error {
	handler, ok := c.Update[lc.LedgerChange.Updated.Data.Type]
	if !ok {
		return errors.From(errors.New("No handler provided for required entry type"), logan.F{
			"entry_type": lc.LedgerChange.Created.Data.Type,
		})
	}

	return handler.Updated(lc)
}

func (c *Consumer) Deleted(lc LedgerChange) error {
	handler, ok := c.Delete[lc.LedgerChange.Removed.Type]
	if !ok {
		return errors.From(errors.New("No handler provided for required entry type"), logan.F{
			"entry_type": lc.LedgerChange.Removed.Type,
		})
	}

	return handler.Deleted(lc)
}
